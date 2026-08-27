package repository

import (
	"kickbase/internal/domain"

	"gorm.io/gorm"
)

type MatchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) Create(match *domain.Match) error {
	return r.db.Create(match).Error
}

func (r *MatchRepository) FindByID(id string) (*domain.Match, error) {
	var match domain.Match
	if err := r.db.First(&match, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) List(page, limit int) ([]domain.Match, int64, error) {
	var matches []domain.Match
	var total int64

	r.db.Model(&domain.Match{}).Count(&total)
	r.db.Offset((page - 1) * limit).Limit(limit).Find(&matches)

	return matches, total, nil
}

func (r *MatchRepository) Update(match *domain.Match) error {
	return r.db.Save(match).Error
}

func (r *MatchRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&domain.Match{}).Where("id = ?", id).Update("status", status).Error
}

func (r *MatchRepository) CreateHistory(history *domain.MatchHistory) error {
	return r.db.Create(history).Error
}

func (r *MatchRepository) GetHistory(matchID string) ([]domain.MatchHistory, error) {
	var history []domain.MatchHistory
	err := r.db.Where("match_id = ?", matchID).Order("version DESC").Find(&history).Error
	return history, err
}

func (r *MatchRepository) GetHistoryByVersion(matchID string, version int) (*domain.MatchHistory, error) {
	var history domain.MatchHistory
	if err := r.db.Where("match_id = ? AND version = ?", matchID, version).First(&history).Error; err != nil {
		return nil, err
	}
	return &history, nil
}
