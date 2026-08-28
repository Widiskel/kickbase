package domain

import (
	"time"
)

type Match struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MatchDate   string    `json:"match_date" gorm:"index;not null;size:10"`
	MatchTime   string    `json:"match_time" gorm:"not null;size:8"`
	HomeTeamID  string    `json:"home_team_id" gorm:"type:uuid;not null;index"`
	AwayTeamID  string    `json:"away_team_id" gorm:"type:uuid;not null;index"`
	Status      string    `json:"status" gorm:"index;not null;size:20;default:'scheduled'"`
	Version     int       `json:"version" gorm:"not null;default:1"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MatchHistory struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MatchID   string    `json:"match_id" gorm:"type:uuid;not null;index"`
	Version   int       `json:"version" gorm:"not null"`
	Changes   string    `json:"changes" gorm:"type:jsonb;not null"`
	CreatedAt time.Time `json:"created_at"`
}
