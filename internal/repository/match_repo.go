package repository

import (
	"kickbase/internal/domain"
	"kickbase/internal/handler"
	"kickbase/internal/interfaces"
	"strings"

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

func (r *MatchRepository) List(opts interfaces.MatchFilterOptions) ([]domain.Match, int64, error) {
	var matches []domain.Match
	var total int64

	query := r.db.Model(&domain.Match{})

	if opts.TeamID != "" {
		fc := handler.ParseFilter(opts.TeamID, handler.OpEQ)
		if fc.Op == handler.OpEQ {
			query = query.Where("home_team_id = ? OR away_team_id = ?", fc.Value, fc.Value)
		} else if fc.Op == handler.OpIN && len(fc.Values) > 0 {
			query = query.Where("home_team_id IN ? OR away_team_id IN ?", fc.Values, fc.Values)
		}
	}
	if opts.Status != "" {
		fc := handler.ParseFilter(opts.Status, handler.OpEQ)
		query = handler.ApplyFilterToQuery(query, "status", fc)
	}
	if opts.DateFrom != "" {
		query = query.Where("match_date >= ?", opts.DateFrom)
	}
	if opts.DateTo != "" {
		query = query.Where("match_date <= ?", opts.DateTo)
	}
	if opts.MatchDate != "" {
		fc := handler.ParseFilter(opts.MatchDate, handler.OpEQ)
		query = handler.ApplyFilterToQuery(query, "match_date", fc)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	sortBy := "match_date"
	order := "ASC"
	if opts.SortBy != "" {
		switch opts.SortBy {
		case "match_date", "match_time", "created_at":
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

	err := query.Offset((opts.Page - 1) * opts.Limit).Limit(opts.Limit).Find(&matches).Error
	return matches, total, err
}

func (r *MatchRepository) CountTotal() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Match{}).Count(&count).Error
	return count, err
}

func (r *MatchRepository) CountCompleted() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Match{}).Where("status = ?", "completed").Count(&count).Error
	return count, err
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
