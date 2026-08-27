package repository

import (
	"kickbase/internal/domain"

	"gorm.io/gorm"
)

type GoalRepository struct {
	db *gorm.DB
}

func NewGoalRepository(db *gorm.DB) *GoalRepository {
	return &GoalRepository{db: db}
}

func (r *GoalRepository) Create(goal *domain.Goal) error {
	return r.db.Create(goal).Error
}

func (r *GoalRepository) FindByID(id string) (*domain.Goal, error) {
	var goal domain.Goal
	if err := r.db.First(&goal, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &goal, nil
}

func (r *GoalRepository) ListByMatchResult(matchResultID string) ([]domain.Goal, error) {
	var goals []domain.Goal
	err := r.db.Where("match_result_id = ?", matchResultID).Find(&goals).Error
	return goals, err
}

func (r *GoalRepository) CountByPlayer(playerID string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Goal{}).Where("player_id = ?", playerID).Count(&count).Error
	return count, err
}

func (r *GoalRepository) CreateHistory(history *domain.GoalHistory) error {
	return r.db.Create(history).Error
}
