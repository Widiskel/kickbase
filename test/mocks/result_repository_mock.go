package mocks

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

var _ interfaces.ResultRepository = (*MockResultRepository)(nil)

type MockResultRepository struct {
	CreateFunc        func(result *domain.MatchResult) error
	FindByMatchIDFunc func(matchID string) (*domain.MatchResult, error)
}

func (m *MockResultRepository) Create(result *domain.MatchResult) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(result)
	}
	return nil
}

func (m *MockResultRepository) FindByMatchID(matchID string) (*domain.MatchResult, error) {
	if m.FindByMatchIDFunc != nil {
		return m.FindByMatchIDFunc(matchID)
	}
	return nil, nil
}
