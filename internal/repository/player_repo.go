package repository

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
	"strings"

	"gorm.io/gorm"
)

type PlayerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Create(player *domain.Player) error {
	return r.db.Create(player).Error
}

func (r *PlayerRepository) FindByID(id string) (*domain.Player, error) {
	var player domain.Player
	if err := r.db.First(&player, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *PlayerRepository) FindByIDIncludingDeleted(id string) (*domain.Player, error) {
	var player domain.Player
	if err := r.db.Unscoped().First(&player, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *PlayerRepository) List(opts interfaces.PlayerFilterOptions) ([]domain.Player, int64, error) {
	var players []domain.Player
	var total int64

	query := r.db.Model(&domain.Player{})

	if opts.TeamID != "" {
		query = query.Where("team_id = ?", opts.TeamID)
	}
	if opts.Name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(opts.Name)+"%")
	}
	if opts.Position != "" {
		query = query.Where("position = ?", opts.Position)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	sortBy := "jersey_number"
	order := "ASC"
	if opts.SortBy != "" {
		switch opts.SortBy {
		case "name", "jersey_number", "height", "weight", "created_at":
			sortBy = opts.SortBy
		}
	}
	if strings.ToUpper(opts.Order) == "DESC" {
		order = "DESC"
	}
	query = query.Order(sortBy + " " + order)

	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	err := query.Offset((opts.Page - 1) * opts.Limit).Limit(opts.Limit).Find(&players).Error
	return players, total, err
}

func (r *PlayerRepository) CountTotal() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Player{}).Count(&count).Error
	return count, err
}

func (r *PlayerRepository) Update(player *domain.Player) error {
	return r.db.Save(player).Error
}

func (r *PlayerRepository) Delete(id string) error {
	return r.db.Delete(&domain.Player{}, "id = ?", id).Error
}

func (r *PlayerRepository) CheckJerseyUnique(teamID string, jerseyNumber int, excludeID string) (bool, error) {
	var count int64
	query := r.db.Model(&domain.Player{}).Where("team_id = ? AND jersey_number = ?", teamID, jerseyNumber)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count == 0, err
}

func (r *PlayerRepository) CountGoals(playerID string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Goal{}).Where("player_id = ?", playerID).Count(&count).Error
	return count, err
}

func (r *PlayerRepository) CreateHistory(history *domain.PlayerHistory) error {
	return r.db.Create(history).Error
}

func (r *PlayerRepository) GetHistory(playerID string) ([]domain.PlayerHistory, error) {
	var history []domain.PlayerHistory
	err := r.db.Where("player_id = ?", playerID).Order("version DESC").Find(&history).Error
	return history, err
}

func (r *PlayerRepository) GetHistoryByVersion(playerID string, version int) (*domain.PlayerHistory, error) {
	var history domain.PlayerHistory
	if err := r.db.Where("player_id = ? AND version = ?", playerID, version).First(&history).Error; err != nil {
		return nil, err
	}
	return &history, nil
}
