package repository

import (
	"kickbase/internal/domain"

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

func (r *TeamRepository) FindByName(name string) (*domain.Team, error) {
	var team domain.Team
	if err := r.db.First(&team, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) List(page, limit int) ([]domain.Team, int64, error) {
	var teams []domain.Team
	var total int64

	r.db.Model(&domain.Team{}).Count(&total)
	r.db.Offset((page - 1) * limit).Limit(limit).Find(&teams)

	return teams, total, nil
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
