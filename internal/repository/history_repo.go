package repository

import (
	"gorm.io/gorm"
)

type HistoryRepository struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (r *HistoryRepository) GetLatestVersion(table string, entityID string) (int, error) {
	var version int
	err := r.db.Table(table).
		Where(getIDColumn(table)+" = ?", entityID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&version).Error
	return version, err
}

func getIDColumn(table string) string {
	switch table {
	case "team_histories":
		return "team_id"
	case "player_histories":
		return "player_id"
	case "match_histories":
		return "match_id"
	case "match_result_histories":
		return "match_result_id"
	case "goal_histories":
		return "goal_id"
	default:
		return "id"
	}
}
