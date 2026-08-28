package interfaces

import "kickbase/internal/domain"

// TeamRepository defines the interface for team data access
type TeamRepository interface {
	Create(team *domain.Team) error
	FindByID(id string) (*domain.Team, error)
	FindByIDIncludingDeleted(id string) (*domain.Team, error)
	FindByName(name string) (*domain.Team, error)
	List(page, limit int) ([]domain.Team, int64, error)
	Update(team *domain.Team) error
	Delete(id string) error
	CountPlayers(teamID string) (int64, error)
	CreateHistory(history *domain.TeamHistory) error
	GetHistory(teamID string) ([]domain.TeamHistory, error)
	GetHistoryByVersion(teamID string, version int) (*domain.TeamHistory, error)
}

// PlayerRepository defines the interface for player data access
type PlayerRepository interface {
	Create(player *domain.Player) error
	FindByID(id string) (*domain.Player, error)
	FindByIDIncludingDeleted(id string) (*domain.Player, error)
	ListByTeam(teamID string, page, limit int) ([]domain.Player, int64, error)
	Update(player *domain.Player) error
	Delete(id string) error
	CheckJerseyUnique(teamID string, jerseyNumber int, excludeID string) (bool, error)
	CountGoals(playerID string) (int64, error)
	CreateHistory(history *domain.PlayerHistory) error
	GetHistory(playerID string) ([]domain.PlayerHistory, error)
	GetHistoryByVersion(playerID string, version int) (*domain.PlayerHistory, error)
}

// MatchRepository defines the interface for match data access
type MatchRepository interface {
	Create(match *domain.Match) error
	FindByID(id string) (*domain.Match, error)
	List(page, limit int) ([]domain.Match, int64, error)
	Update(match *domain.Match) error
	CreateHistory(history *domain.MatchHistory) error
	GetHistory(matchID string) ([]domain.MatchHistory, error)
	GetHistoryByVersion(matchID string, version int) (*domain.MatchHistory, error)
}

// ResultRepository defines the interface for match result data access
type ResultRepository interface {
	Create(result *domain.MatchResult) error
	FindByMatchID(matchID string) (*domain.MatchResult, error)
}

// GoalRepository defines the interface for goal data access
type GoalRepository interface {
	Create(goal *domain.Goal) error
	ListByMatchResult(resultID string) ([]domain.Goal, error)
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *domain.User) error
	FindByUsername(username string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
}

// RefreshTokenRepository defines the interface for refresh token data access
type RefreshTokenRepository interface {
	Create(token *domain.RefreshToken) error
	FindByToken(token string) (*domain.RefreshToken, error)
	Revoke(token string) error
	RevokeAllForUser(userID string) error
}
