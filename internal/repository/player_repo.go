package repository

import (
	"kickbase/internal/domain"

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

func (r *PlayerRepository) ListByTeam(teamID string, page, limit int) ([]domain.Player, int64, error) {
	var players []domain.Player
	var total int64

	r.db.Model(&domain.Player{}).Where("team_id = ?", teamID).Count(&total)
	r.db.Where("team_id = ?", teamID).Offset((page - 1) * limit).Limit(limit).Find(&players)

	return players, total, nil
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
