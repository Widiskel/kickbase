package mocks

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

var _ interfaces.GoalRepository = (*MockGoalRepository)(nil)

type MockGoalRepository struct {
	CreateFunc            func(goal *domain.Goal) error
	ListByMatchResultFunc func(resultID string) ([]domain.Goal, error)
}

func (m *MockGoalRepository) Create(goal *domain.Goal) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(goal)
	}
	return nil
}

func (m *MockGoalRepository) ListByMatchResult(resultID string) ([]domain.Goal, error) {
	if m.ListByMatchResultFunc != nil {
		return m.ListByMatchResultFunc(resultID)
	}
	return nil, nil
}
