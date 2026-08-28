package service_test

import (
	"errors"
	"testing"

	"kickbase/internal/domain"
	"kickbase/internal/service"
	"kickbase/test/mocks"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return nil, nil // No duplicate
		},
		CreateFunc: func(user *domain.User) error {
			user.ID = "user-1"
			return nil
		},
	}

	svc := service.NewAuthService(mockRepo, "test-secret")

	user, err := svc.Register("admin_xyz", "password123", "Admin XYZ", "admin")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "admin_xyz", user.Username)
	assert.Equal(t, "admin", user.Role)
	assert.NotEmpty(t, user.Password) // Hashed
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return &domain.User{Username: username}, nil
		},
	}

	svc := service.NewAuthService(mockRepo, "test-secret")

	user, err := svc.Register("admin_xyz", "password123", "Admin XYZ", "admin")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "username already exists", err.Error())
}

func TestAuthService_Login_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo := &mocks.MockUserRepository{
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

	svc := service.NewAuthService(mockRepo, "test-secret")

	token, user, err := svc.Login("admin_xyz", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotNil(t, user)
	assert.Equal(t, "admin_xyz", user.Username)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo := &mocks.MockUserRepository{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return &domain.User{
				ID:       "user-1",
				Username: username,
				Password: string(hashedPassword),
				Role:     "admin",
			}, nil
		},
	}

	svc := service.NewAuthService(mockRepo, "test-secret")

	token, user, err := svc.Login("admin_xyz", "wrong-password")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	assert.Equal(t, "invalid username or password", err.Error())
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return nil, errors.New("user not found")
		},
	}

	svc := service.NewAuthService(mockRepo, "test-secret")

	token, user, err := svc.Login("non-existent", "password123")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	assert.Equal(t, "invalid username or password", err.Error())
}
