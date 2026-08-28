package interfaces

import "kickbase/internal/domain"

// TeamService defines the interface for team business logic
type TeamService interface {
	CreateTeam(team *domain.Team) error
	GetTeam(id string) (*domain.Team, error)
	ListTeams(page, limit int) ([]domain.Team, int64, error)
	UpdateTeam(team *domain.Team) error
	DeleteTeam(id string) error
	GetTeamHistory(teamID string) ([]domain.TeamHistory, error)
	RevertTeam(teamID string, targetVersion int) error
}

// PlayerService defines the interface for player business logic
type PlayerService interface {
	CreatePlayer(player *domain.Player) error
	GetPlayer(id string) (*domain.Player, error)
	ListPlayersByTeam(teamID string, page, limit int) ([]domain.Player, int64, error)
	UpdatePlayer(player *domain.Player) error
	DeletePlayer(id string) error
	GetPlayerHistory(playerID string) ([]domain.PlayerHistory, error)
	RevertPlayer(playerID string, targetVersion int) error
}

// MatchService defines the interface for match business logic
type MatchService interface {
	CreateMatch(match *domain.Match) error
	GetMatch(id string) (*domain.Match, error)
	ListMatches(page, limit int) ([]domain.Match, int64, error)
	UpdateMatchStatus(id string, status string) error
	GetMatchHistory(matchID string) ([]domain.MatchHistory, error)
	RevertMatch(matchID string, targetVersion int) error
}

// ResultService defines the interface for match result business logic
type ResultService interface {
	CreateResult(input CreateResultInput) (*domain.MatchResult, error)
	GetResult(matchID string) (*domain.MatchResult, []domain.Goal, error)
}

// ReportService defines the interface for report business logic
type ReportService interface {
	GetMatchReport(matchID string) (*MatchReport, error)
	ListMatchReports(page, limit int) ([]MatchReport, int64, error)
}

// AuthService defines the interface for authentication & authorization business logic
type AuthService interface {
	Register(username, password, name, role string) (*domain.User, error)
	Login(username, password string) (string, *domain.User, error)
}

// CreateResultInput represents the input for creating a match result
type CreateResultInput struct {
	MatchID   string      `json:"match_id"`
	HomeScore int         `json:"home_score"`
	AwayScore int         `json:"away_score"`
	Goals     []GoalInput `json:"goals"`
}

// GoalInput represents a goal in the result input
type GoalInput struct {
	PlayerID string `json:"player_id"`
	GoalTime string `json:"goal_time"`
}

// MatchReport represents a match report
type MatchReport struct {
	MatchID            string      `json:"match_id"`
	MatchDate          string      `json:"match_date"`
	MatchTime          string      `json:"match_time"`
	HomeTeam           string      `json:"home_team"`
	AwayTeam           string      `json:"away_team"`
	HomeScore          *int        `json:"home_score,omitempty"`
	AwayScore          *int        `json:"away_score,omitempty"`
	Status             string      `json:"status"`
	TopScorers         []TopScorer `json:"top_scorers,omitempty"`
	CumulativeHomeWins int         `json:"cumulative_home_wins"`
	CumulativeAwayWins int         `json:"cumulative_away_wins"`
}

// TopScorer represents a top goal scorer
type TopScorer struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Goals      int    `json:"goals"`
}
