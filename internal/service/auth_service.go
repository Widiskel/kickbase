package service

import (
	"errors"
	"fmt"
	"time"

	"kickbase/internal/domain"
	"kickbase/internal/interfaces"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  interfaces.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo interfaces.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(username, password, name, role string) (*domain.User, error) {
	if username == "" || password == "" || name == "" {
		return nil, errors.New("username, password, and name are required")
	}

	if role == "" {
		role = "admin"
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

func (s *AuthService) Login(username, password string) (string, *domain.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil || user == nil {
		return "", nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, user, nil
}
