package mocks

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

var _ interfaces.RefreshTokenRepository = (*MockRefreshTokenRepository)(nil)

type MockRefreshTokenRepository struct {
	CreateFunc           func(token *domain.RefreshToken) error
	FindByTokenFunc      func(token string) (*domain.RefreshToken, error)
	RevokeFunc           func(token string) error
	RevokeAllForUserFunc func(userID string) error
}

func (m *MockRefreshTokenRepository) Create(token *domain.RefreshToken) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(token)
	}
	return nil
}

func (m *MockRefreshTokenRepository) FindByToken(token string) (*domain.RefreshToken, error) {
	if m.FindByTokenFunc != nil {
		return m.FindByTokenFunc(token)
	}
	return nil, nil
}

func (m *MockRefreshTokenRepository) Revoke(token string) error {
	if m.RevokeFunc != nil {
		return m.RevokeFunc(token)
	}
	return nil
}

func (m *MockRefreshTokenRepository) RevokeAllForUser(userID string) error {
	if m.RevokeAllForUserFunc != nil {
		return m.RevokeAllForUserFunc(userID)
	}
	return nil
}
