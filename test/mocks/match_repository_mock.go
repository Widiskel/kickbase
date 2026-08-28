package mocks

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

var _ interfaces.MatchRepository = (*MockMatchRepository)(nil)

type MockMatchRepository struct {
	CreateFunc              func(match *domain.Match) error
	FindByIDFunc            func(id string) (*domain.Match, error)
	ListFunc                func(opts interfaces.MatchFilterOptions) ([]domain.Match, int64, error)
	CountTotalFunc          func() (int64, error)
	CountCompletedFunc      func() (int64, error)
	UpdateFunc              func(match *domain.Match) error
	CreateHistoryFunc       func(history *domain.MatchHistory) error
	GetHistoryFunc          func(matchID string) ([]domain.MatchHistory, error)
	GetHistoryByVersionFunc func(matchID string, version int) (*domain.MatchHistory, error)
}

func (m *MockMatchRepository) Create(match *domain.Match) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(match)
	}
	return nil
}

func (m *MockMatchRepository) FindByID(id string) (*domain.Match, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockMatchRepository) List(opts interfaces.MatchFilterOptions) ([]domain.Match, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(opts)
	}
	return nil, 0, nil
}

func (m *MockMatchRepository) CountTotal() (int64, error) {
	if m.CountTotalFunc != nil {
		return m.CountTotalFunc()
	}
	return 0, nil
}

func (m *MockMatchRepository) CountCompleted() (int64, error) {
	if m.CountCompletedFunc != nil {
		return m.CountCompletedFunc()
	}
	return 0, nil
}

func (m *MockMatchRepository) Update(match *domain.Match) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(match)
	}
	return nil
}

func (m *MockMatchRepository) CreateHistory(history *domain.MatchHistory) error {
	if m.CreateHistoryFunc != nil {
		return m.CreateHistoryFunc(history)
	}
	return nil
}

func (m *MockMatchRepository) GetHistory(matchID string) ([]domain.MatchHistory, error) {
	if m.GetHistoryFunc != nil {
		return m.GetHistoryFunc(matchID)
	}
	return nil, nil
}

func (m *MockMatchRepository) GetHistoryByVersion(matchID string, version int) (*domain.MatchHistory, error) {
	if m.GetHistoryByVersionFunc != nil {
		return m.GetHistoryByVersionFunc(matchID, version)
	}
	return nil, nil
}
