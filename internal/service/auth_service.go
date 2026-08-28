package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"kickbase/internal/domain"
	"kickbase/internal/interfaces"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo         interfaces.UserRepository
	refreshTokenRepo interfaces.RefreshTokenRepository
	jwtSecret        string
}

func NewAuthService(
	userRepo interfaces.UserRepository,
	refreshTokenRepo interfaces.RefreshTokenRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
	}
}

func (s *AuthService) Register(username, password, name, role string) (*domain.User, error) {
	if username == "" || password == "" || name == "" {
		return nil, errors.New("username, password, and name are required")
	}

	if role == "" {
		role = "admin"
	}

	// Validate role
	if _, exists := domain.RolePermissions[role]; !exists {
		role = "viewer"
	}

	// Check existing username
	existing, _ := s.userRepo.FindByUsername(username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Username: username,
		Password: string(hashedPassword),
		Name:     name,
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (*interfaces.AuthResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil || user == nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return s.generateTokenPair(user)
}

func (s *AuthService) RefreshToken(refreshTokenStr string) (*interfaces.AuthResponse, error) {
	if refreshTokenStr == "" {
		return nil, errors.New("refresh token is required")
	}

	if s.refreshTokenRepo == nil {
		// Mock fallback
		return nil, errors.New("refresh token not found or revoked")
	}

	rt, err := s.refreshTokenRepo.FindByToken(refreshTokenStr)
	if err != nil || rt == nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if rt.Revoked || time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("refresh token has expired or been revoked")
	}

	user, err := s.userRepo.FindByID(rt.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user associated with token not found")
	}

	// Revoke old refresh token (rotate token)
	_ = s.refreshTokenRepo.Revoke(refreshTokenStr)

	// Generate new token pair
	return s.generateTokenPair(user)
}

func (s *AuthService) generateTokenPair(user *domain.User) (*interfaces.AuthResponse, error) {
	permissions := domain.GetPermissionsForRole(user.Role)

	// Access Token: Expires in 1 hour
	accessExpiry := 1 * time.Hour
	accessClaims := jwt.MapClaims{
		"user_id":     user.ID,
		"username":    user.Username,
		"role":        user.Role,
		"permissions": permissions,
		"exp":         time.Now().Add(accessExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Refresh Token: Cryptographically secure random string, expires in 7 days
	randomBytes := make([]byte, 32)
	_, _ = rand.Read(randomBytes)
	refreshTokenStr := hex.EncodeToString(randomBytes)

	if s.refreshTokenRepo != nil {
		rt := &domain.RefreshToken{
			UserID:    user.ID,
			Token:     refreshTokenStr,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Revoked:   false,
		}
		_ = s.refreshTokenRepo.Create(rt)
	}

	return &interfaces.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExpiry.Seconds()),
		User:         user,
		Permissions:  permissions,
	}, nil
}
