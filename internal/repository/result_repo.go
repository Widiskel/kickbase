package repository

import (
	"kickbase/internal/domain"

	"gorm.io/gorm"
)

type ResultRepository struct {
	db *gorm.DB
}

func NewResultRepository(db *gorm.DB) *ResultRepository {
	return &ResultRepository{db: db}
}

func (r *ResultRepository) Create(result *domain.MatchResult) error {
	return r.db.Create(result).Error
}

func (r *ResultRepository) FindByID(id string) (*domain.MatchResult, error) {
	var result domain.MatchResult
	if err := r.db.First(&result, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ResultRepository) FindByMatchID(matchID string) (*domain.MatchResult, error) {
	var result domain.MatchResult
	if err := r.db.Where("match_id = ?", matchID).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ResultRepository) Update(result *domain.MatchResult) error {
	return r.db.Save(result).Error
}

func (r *ResultRepository) CreateHistory(history *domain.MatchResultHistory) error {
	return r.db.Create(history).Error
}

func (r *ResultRepository) GetHistory(resultID string) ([]domain.MatchResultHistory, error) {
	var history []domain.MatchResultHistory
	err := r.db.Where("match_result_id = ?", resultID).Order("version DESC").Find(&history).Error
	return history, err
}
