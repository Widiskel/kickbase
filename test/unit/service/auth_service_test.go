package service_test

import (
	"errors"
	"testing"
	"time"

	"kickbase/internal/domain"
	"kickbase/internal/service"
	"kickbase/test/mocks"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Success(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return nil, nil // No duplicate
		},
		CreateFunc: func(user *domain.User) error {
			user.ID = "user-1"
			return nil
		},
	}
	mockTokenRepo := &mocks.MockRefreshTokenRepository{}

	svc := service.NewAuthService(mockUserRepo, mockTokenRepo, "test-secret")

	user, err := svc.Register("admin_xyz", "password123", "Admin XYZ", "admin")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "admin_xyz", user.Username)
	assert.Equal(t, "admin", user.Role)
	assert.NotEmpty(t, user.Password) // Hashed
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return &domain.User{Username: username}, nil
		},
	}

	svc := service.NewAuthService(mockUserRepo, nil, "test-secret")

	user, err := svc.Register("admin_xyz", "password123", "Admin XYZ", "admin")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "username already exists", err.Error())
}

func TestAuthService_Login_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUserRepo := &mocks.MockUserRepository{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return &domain.User{
				ID:       "user-1",
				Username: username,
				Password: string(hashedPassword),
				Name:     "Admin XYZ",
				Role:     "admin",
			}, nil
		},
	}
	mockTokenRepo := &mocks.MockRefreshTokenRepository{
		CreateFunc: func(token *domain.RefreshToken) error {
			return nil
		},
	}

	svc := service.NewAuthService(mockUserRepo, mockTokenRepo, "test-secret")

	authResp, err := svc.Login("admin_xyz", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, authResp)
	assert.NotEmpty(t, authResp.AccessToken)
	assert.NotEmpty(t, authResp.RefreshToken)
	assert.Equal(t, "Bearer", authResp.TokenType)
	assert.Contains(t, authResp.Permissions, domain.PermTeamsCreate)
	assert.Contains(t, authResp.Permissions, domain.PermTeamsDelete)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUserRepo := &mocks.MockUserRepository{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return &domain.User{
				ID:       "user-1",
				Username: username,
				Password: string(hashedPassword),
				Role:     "admin",
			}, nil
		},
	}

	svc := service.NewAuthService(mockUserRepo, nil, "test-secret")

	authResp, err := svc.Login("admin_xyz", "wrong-password")
	assert.Error(t, err)
	assert.Nil(t, authResp)
	assert.Equal(t, "invalid username or password", err.Error())
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return nil, errors.New("user not found")
		},
	}

	svc := service.NewAuthService(mockUserRepo, nil, "test-secret")

	authResp, err := svc.Login("non-existent", "password123")
	assert.Error(t, err)
	assert.Nil(t, authResp)
	assert.Equal(t, "invalid username or password", err.Error())
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{
		FindByIDFunc: func(id string) (*domain.User, error) {
			return &domain.User{
				ID:       "user-1",
				Username: "admin_xyz",
				Role:     "admin",
			}, nil
		},
	}
	mockTokenRepo := &mocks.MockRefreshTokenRepository{
		FindByTokenFunc: func(token string) (*domain.RefreshToken, error) {
			return &domain.RefreshToken{
				UserID:    "user-1",
				Token:     token,
				ExpiresAt: time.Now().Add(24 * time.Hour),
				Revoked:   false,
			}, nil
		},
		RevokeFunc: func(token string) error {
			return nil
		},
		CreateFunc: func(token *domain.RefreshToken) error {
			return nil
		},
	}

	svc := service.NewAuthService(mockUserRepo, mockTokenRepo, "test-secret")

	authResp, err := svc.RefreshToken("valid-refresh-token")
	assert.NoError(t, err)
	assert.NotNil(t, authResp)
	assert.NotEmpty(t, authResp.AccessToken)
	assert.NotEmpty(t, authResp.RefreshToken)
}

func TestAuthService_RefreshToken_ExpiredOrRevoked(t *testing.T) {
	mockTokenRepo := &mocks.MockRefreshTokenRepository{
		FindByTokenFunc: func(token string) (*domain.RefreshToken, error) {
			return &domain.RefreshToken{
				UserID:    "user-1",
				Token:     token,
				ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
				Revoked:   false,
			}, nil
		},
	}

	svc := service.NewAuthService(nil, mockTokenRepo, "test-secret")

	authResp, err := svc.RefreshToken("expired-token")
	assert.Error(t, err)
	assert.Nil(t, authResp)
	assert.Equal(t, "refresh token has expired or been revoked", err.Error())
}

func TestDomain_Permissions(t *testing.T) {
	// Admin has all permissions
	assert.True(t, domain.HasPermission("admin", domain.PermTeamsCreate))
	assert.True(t, domain.HasPermission("admin", domain.PermTeamsDelete))
	assert.True(t, domain.HasPermission("admin", domain.PermTeamsRevert))

	// Staff has create & update but NOT delete or revert
	assert.True(t, domain.HasPermission("staff", domain.PermTeamsCreate))
	assert.True(t, domain.HasPermission("staff", domain.PermTeamsUpdate))
	assert.False(t, domain.HasPermission("staff", domain.PermTeamsDelete))
	assert.False(t, domain.HasPermission("staff", domain.PermTeamsRevert))

	// Viewer has ONLY read permissions
	assert.True(t, domain.HasPermission("viewer", domain.PermTeamsRead))
	assert.False(t, domain.HasPermission("viewer", domain.PermTeamsCreate))
	assert.False(t, domain.HasPermission("viewer", domain.PermTeamsDelete))
}
