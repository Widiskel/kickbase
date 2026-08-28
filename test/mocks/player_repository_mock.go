package mocks

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

var _ interfaces.PlayerRepository = (*MockPlayerRepository)(nil)

type MockPlayerRepository struct {
	CreateFunc                   func(player *domain.Player) error
	FindByIDFunc                 func(id string) (*domain.Player, error)
	FindByIDIncludingDeletedFunc func(id string) (*domain.Player, error)
	ListFunc                     func(opts interfaces.PlayerFilterOptions) ([]domain.Player, int64, error)
	CountTotalFunc               func() (int64, error)
	UpdateFunc                   func(player *domain.Player) error
	DeleteFunc                   func(id string) error
	CheckJerseyUniqueFunc        func(teamID string, jerseyNumber int, excludeID string) (bool, error)
	CountGoalsFunc               func(playerID string) (int64, error)
	CreateHistoryFunc            func(history *domain.PlayerHistory) error
	GetHistoryFunc               func(playerID string) ([]domain.PlayerHistory, error)
	GetHistoryByVersionFunc      func(playerID string, version int) (*domain.PlayerHistory, error)
}

func (m *MockPlayerRepository) Create(player *domain.Player) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(player)
	}
	return nil
}

func (m *MockPlayerRepository) FindByID(id string) (*domain.Player, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockPlayerRepository) FindByIDIncludingDeleted(id string) (*domain.Player, error) {
	if m.FindByIDIncludingDeletedFunc != nil {
		return m.FindByIDIncludingDeletedFunc(id)
	}
	return nil, nil
}

func (m *MockPlayerRepository) List(opts interfaces.PlayerFilterOptions) ([]domain.Player, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(opts)
	}
	return nil, 0, nil
}

func (m *MockPlayerRepository) CountTotal() (int64, error) {
	if m.CountTotalFunc != nil {
		return m.CountTotalFunc()
	}
	return 0, nil
}

func (m *MockPlayerRepository) Update(player *domain.Player) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(player)
	}
	return nil
}

func (m *MockPlayerRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockPlayerRepository) CheckJerseyUnique(teamID string, jerseyNumber int, excludeID string) (bool, error) {
	if m.CheckJerseyUniqueFunc != nil {
		return m.CheckJerseyUniqueFunc(teamID, jerseyNumber, excludeID)
	}
	return true, nil
}

func (m *MockPlayerRepository) CountGoals(playerID string) (int64, error) {
	if m.CountGoalsFunc != nil {
		return m.CountGoalsFunc(playerID)
	}
	return 0, nil
}

func (m *MockPlayerRepository) CreateHistory(history *domain.PlayerHistory) error {
	if m.CreateHistoryFunc != nil {
		return m.CreateHistoryFunc(history)
	}
	return nil
}

func (m *MockPlayerRepository) GetHistory(playerID string) ([]domain.PlayerHistory, error) {
	if m.GetHistoryFunc != nil {
		return m.GetHistoryFunc(playerID)
	}
	return nil, nil
}

func (m *MockPlayerRepository) GetHistoryByVersion(playerID string, version int) (*domain.PlayerHistory, error) {
	if m.GetHistoryByVersionFunc != nil {
		return m.GetHistoryByVersionFunc(playerID, version)
	}
	return nil, nil
}
