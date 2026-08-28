package handler

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

// Team DTOs

type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	LogoURL     string `json:"logo_url"`
	FoundedYear int    `json:"founded_year" binding:"required,min=1"`
	Address     string `json:"address" binding:"required"`
	City        string `json:"city" binding:"required"`
}

func (r *CreateTeamRequest) ToDomain() *domain.Team {
	return &domain.Team{
		Name:        r.Name,
		LogoURL:     r.LogoURL,
		FoundedYear: r.FoundedYear,
		Address:     r.Address,
		City:        r.City,
	}
}

type UpdateTeamRequest struct {
	Name        string `json:"name"`
	LogoURL     string `json:"logo_url"`
	FoundedYear int    `json:"founded_year"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Version     int    `json:"version" binding:"required"`
}

func (r *UpdateTeamRequest) ToDomain(id string) *domain.Team {
	return &domain.Team{
		ID:          id,
		Name:        r.Name,
		LogoURL:     r.LogoURL,
		FoundedYear: r.FoundedYear,
		Address:     r.Address,
		City:        r.City,
		Version:     r.Version,
	}
}

// Player DTOs

type CreatePlayerRequest struct {
	TeamID       string  `json:"team_id" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Height       float64 `json:"height" binding:"required,gt=0"`
	Weight       float64 `json:"weight" binding:"required,gt=0"`
	Position     string  `json:"position" binding:"required"`
	Playstyle    *string `json:"playstyle"`
	JerseyNumber int     `json:"jersey_number" binding:"required,gt=0"`
}

func (r *CreatePlayerRequest) ToDomain() *domain.Player {
	return &domain.Player{
		TeamID:       r.TeamID,
		Name:         r.Name,
		Height:       r.Height,
		Weight:       r.Weight,
		Position:     r.Position,
		Playstyle:    r.Playstyle,
		JerseyNumber: r.JerseyNumber,
	}
}

type UpdatePlayerRequest struct {
	Name         string  `json:"name" binding:"required"`
	Height       float64 `json:"height" binding:"required,gt=0"`
	Weight       float64 `json:"weight" binding:"required,gt=0"`
	Position     string  `json:"position" binding:"required"`
	Playstyle    *string `json:"playstyle"`
	JerseyNumber int     `json:"jersey_number" binding:"required,gt=0"`
	Version      int     `json:"version" binding:"required"`
}

func (r *UpdatePlayerRequest) ToDomain(id string) *domain.Player {
	return &domain.Player{
		ID:           id,
		Name:         r.Name,
		Height:       r.Height,
		Weight:       r.Weight,
		Position:     r.Position,
		Playstyle:    r.Playstyle,
		JerseyNumber: r.JerseyNumber,
		Version:      r.Version,
	}
}

// Match DTOs

type CreateMatchRequest struct {
	MatchDate  string `json:"match_date" binding:"required"`
	MatchTime  string `json:"match_time" binding:"required"`
	HomeTeamID string `json:"home_team_id" binding:"required"`
	AwayTeamID string `json:"away_team_id" binding:"required"`
}

func (r *CreateMatchRequest) ToDomain() *domain.Match {
	return &domain.Match{
		MatchDate:  r.MatchDate,
		MatchTime:  r.MatchTime,
		HomeTeamID: r.HomeTeamID,
		AwayTeamID: r.AwayTeamID,
	}
}

type UpdateMatchStatusRequest struct {
	Status  string `json:"status" binding:"required"`
	Version int    `json:"version" binding:"required"`
}

// Result DTOs

type CreateResultRequest struct {
	MatchID   string     `json:"match_id" binding:"required"`
	HomeScore int        `json:"home_score" binding:"min=0"`
	AwayScore int        `json:"away_score" binding:"min=0"`
	Goals     []GoalInput `json:"goals"`
}

type GoalInput struct {
	PlayerID string `json:"player_id" binding:"required"`
	GoalTime string `json:"goal_time" binding:"required"`
}

func (r *CreateResultRequest) ToServiceInput() interfaces.CreateResultInput {
	goals := make([]interfaces.GoalInput, len(r.Goals))
	for i, g := range r.Goals {
		goals[i] = interfaces.GoalInput{
			PlayerID: g.PlayerID,
			GoalTime: g.GoalTime,
		}
	}
	return interfaces.CreateResultInput{
		MatchID:   r.MatchID,
		HomeScore: r.HomeScore,
		AwayScore: r.AwayScore,
		Goals:     goals,
	}
}

// Common DTOs

type RevertRequest struct {
	TargetVersion int `json:"target_version" binding:"required"`
}

// Auth DTOs

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
