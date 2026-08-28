package repository

import (
	"kickbase/internal/domain"
	"kickbase/internal/handler"
	"kickbase/internal/interfaces"
	"strings"

	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(team *domain.Team) error {
	return r.db.Create(team).Error
}

func (r *TeamRepository) FindByID(id string) (*domain.Team, error) {
	var team domain.Team
	if err := r.db.First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) FindByIDIncludingDeleted(id string) (*domain.Team, error) {
	var team domain.Team
	if err := r.db.Unscoped().First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) FindByName(name string) (*domain.Team, error) {
	var team domain.Team
	if err := r.db.First(&team, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) List(opts interfaces.TeamFilterOptions) ([]domain.Team, int64, error) {
	var teams []domain.Team
	var total int64

	query := r.db.Model(&domain.Team{})

	if opts.Name != "" {
		fc := handler.ParseFilter(opts.Name, handler.OpCT)
		query = handler.ApplyFilterToQuery(query, "name", fc)
	}
	if opts.City != "" {
		fc := handler.ParseFilter(opts.City, handler.OpEQ)
		query = handler.ApplyFilterToQuery(query, "city", fc)
	}
	if opts.FoundedYear != "" {
		fc := handler.ParseFilter(opts.FoundedYear, handler.OpEQ)
		query = handler.ApplyFilterToQuery(query, "founded_year", fc)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	sortBy := "created_at"
	order := "ASC"
	if opts.SortBy != "" {
		switch opts.SortBy {
		case "name", "city", "founded_year", "created_at":
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

	err := query.Offset((opts.Page - 1) * opts.Limit).Limit(opts.Limit).Find(&teams).Error
	return teams, total, err
}

func (r *TeamRepository) CountTotal() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Team{}).Count(&count).Error
	return count, err
}

func (r *TeamRepository) Update(team *domain.Team) error {
	return r.db.Save(team).Error
}

func (r *TeamRepository) Delete(id string) error {
	return r.db.Delete(&domain.Team{}, "id = ?", id).Error
}

func (r *TeamRepository) CountPlayers(teamID string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Player{}).Where("team_id = ?", teamID).Count(&count).Error
	return count, err
}

func (r *TeamRepository) CreateHistory(history *domain.TeamHistory) error {
	return r.db.Create(history).Error
}

func (r *TeamRepository) GetHistory(teamID string) ([]domain.TeamHistory, error) {
	var history []domain.TeamHistory
	err := r.db.Where("team_id = ?", teamID).Order("version DESC").Find(&history).Error
	return history, err
}

func (r *TeamRepository) GetHistoryByVersion(teamID string, version int) (*domain.TeamHistory, error) {
	var history domain.TeamHistory
	if err := r.db.Where("team_id = ? AND version = ?", teamID, version).First(&history).Error; err != nil {
		return nil, err
	}
	return &history, nil
}
