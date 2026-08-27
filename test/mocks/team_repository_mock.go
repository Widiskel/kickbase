package mocks

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

// Ensure MockTeamRepository implements interfaces.TeamRepository
var _ interfaces.TeamRepository = (*MockTeamRepository)(nil)

// MockTeamRepository is a mock implementation of TeamRepository
type MockTeamRepository struct {
	CreateFunc                 func(team *domain.Team) error
	FindByIDFunc               func(id string) (*domain.Team, error)
	FindByIDIncludingDeletedFunc func(id string) (*domain.Team, error)
	FindByNameFunc             func(name string) (*domain.Team, error)
	ListFunc                   func(page, limit int) ([]domain.Team, int64, error)
	UpdateFunc                 func(team *domain.Team) error
	DeleteFunc                 func(id string) error
	CountPlayersFunc           func(teamID string) (int64, error)
	CreateHistoryFunc          func(history *domain.TeamHistory) error
	GetHistoryFunc             func(teamID string) ([]domain.TeamHistory, error)
	GetHistoryByVersionFunc    func(teamID string, version int) (*domain.TeamHistory, error)
}

func (m *MockTeamRepository) Create(team *domain.Team) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(team)
	}
	return nil
}

func (m *MockTeamRepository) FindByID(id string) (*domain.Team, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockTeamRepository) FindByIDIncludingDeleted(id string) (*domain.Team, error) {
	if m.FindByIDIncludingDeletedFunc != nil {
		return m.FindByIDIncludingDeletedFunc(id)
	}
	return nil, nil
}

func (m *MockTeamRepository) FindByName(name string) (*domain.Team, error) {
	if m.FindByNameFunc != nil {
		return m.FindByNameFunc(name)
	}
	return nil, nil
}

func (m *MockTeamRepository) List(page, limit int) ([]domain.Team, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(page, limit)
	}
	return nil, 0, nil
}

func (m *MockTeamRepository) Update(team *domain.Team) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(team)
	}
	return nil
}

func (m *MockTeamRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockTeamRepository) CountPlayers(teamID string) (int64, error) {
	if m.CountPlayersFunc != nil {
		return m.CountPlayersFunc(teamID)
	}
	return 0, nil
}

func (m *MockTeamRepository) CreateHistory(history *domain.TeamHistory) error {
	if m.CreateHistoryFunc != nil {
		return m.CreateHistoryFunc(history)
	}
	return nil
}

func (m *MockTeamRepository) GetHistory(teamID string) ([]domain.TeamHistory, error) {
	if m.GetHistoryFunc != nil {
		return m.GetHistoryFunc(teamID)
	}
	return nil, nil
}

func (m *MockTeamRepository) GetHistoryByVersion(teamID string, version int) (*domain.TeamHistory, error) {
	if m.GetHistoryByVersionFunc != nil {
		return m.GetHistoryByVersionFunc(teamID, version)
	}
	return nil, nil
}
