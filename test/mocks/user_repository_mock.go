package mocks

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

var _ interfaces.UserRepository = (*MockUserRepository)(nil)

type MockUserRepository struct {
	CreateFunc         func(user *domain.User) error
	FindByUsernameFunc func(username string) (*domain.User, error)
	FindByIDFunc       func(id string) (*domain.User, error)
}

func (m *MockUserRepository) Create(user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

func (m *MockUserRepository) FindByUsername(username string) (*domain.User, error) {
	if m.FindByUsernameFunc != nil {
		return m.FindByUsernameFunc(username)
	}
	return nil, nil
}

func (m *MockUserRepository) FindByID(id string) (*domain.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
