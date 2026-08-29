package handler

import (
	"kickbase/internal/domain"
	"kickbase/internal/interfaces"
)

// Team DTOs

type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required" example:"Persija Jakarta"`
	LogoURL     string `json:"logo_url" example:"https://upload.wikimedia.org/wikipedia/id/5/5e/Logo_Persija.png"`
	FoundedYear int    `json:"founded_year" binding:"required,min=1" example:"1928"`
	Address     string `json:"address" binding:"required" example:"Jl. Rasuna Said Kav. C-22"`
	City        string `json:"city" binding:"required" example:"Jakarta"`
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
	Name        string `json:"name" example:"Persija Jakarta Updated"`
	LogoURL     string `json:"logo_url" example:"https://upload.wikimedia.org/wikipedia/id/5/5e/Logo_Persija.png"`
	FoundedYear int    `json:"founded_year" example:"1928"`
	Address     string `json:"address" example:"Jl. Rasuna Said Kav. C-22 Baru"`
	City        string `json:"city" example:"Jakarta"`
	Version     int    `json:"version" binding:"required" example:"1"`
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
	TeamID       string  `json:"team_id" binding:"required" example:"00000000-0000-0000-0000-000000000001"`
	Name         string  `json:"name" binding:"required" example:"Marko Simic"`
	Height       float64 `json:"height" binding:"required,gt=0" example:"187"`
	Weight       float64 `json:"weight" binding:"required,gt=0" example:"84"`
	Position     string  `json:"position" binding:"required" example:"CF"`
	Playstyle    *string `json:"playstyle" example:"Goal Poacher"`
	JerseyNumber int     `json:"jersey_number" binding:"required,gt=0" example:"9"`
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
	Name         string  `json:"name" binding:"required" example:"Marko Simic"`
	Height       float64 `json:"height" binding:"required,gt=0" example:"187"`
	Weight       float64 `json:"weight" binding:"required,gt=0" example:"84"`
	Position     string  `json:"position" binding:"required" example:"CF"`
	Playstyle    *string `json:"playstyle" example:"Goal Poacher"`
	JerseyNumber int     `json:"jersey_number" binding:"required,gt=0" example:"9"`
	Version      int     `json:"version" binding:"required" example:"1"`
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
	MatchDate  string `json:"match_date" binding:"required" example:"2026-09-15"`
	MatchTime  string `json:"match_time" binding:"required" example:"19:30:00"`
	HomeTeamID string `json:"home_team_id" binding:"required" example:"00000000-0000-0000-0000-000000000001"`
	AwayTeamID string `json:"away_team_id" binding:"required" example:"00000000-0000-0000-0000-000000000002"`
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
	Status  string `json:"status" binding:"required" example:"completed"`
	Version int    `json:"version" binding:"required" example:"1"`
}

// Result DTOs

type CreateResultRequest struct {
	MatchID   string      `json:"match_id" binding:"required" example:"00000000-0000-0000-0000-000000000001"`
	HomeScore int         `json:"home_score" binding:"min=0" example:"2"`
	AwayScore int         `json:"away_score" binding:"min=0" example:"1"`
	Goals     []GoalInput `json:"goals" binding:"required"`
}

type GoalInput struct {
	PlayerID string `json:"player_id" binding:"required" example:"00000000-0000-0000-0000-000000000001"`
	GoalTime string `json:"goal_time" binding:"required" example:"24"`
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

type RevertRequest struct {
	TargetVersion int `json:"target_version" binding:"required" example:"1"`
}

// Auth DTOs

type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"admin_staff"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Name     string `json:"name" binding:"required" example:"Admin Staff"`
	Role     string `json:"role" binding:"required,oneof=admin staff viewer" example:"staff"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`
}
