package domain

import (
	"time"
)

type MatchResult struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MatchID    string    `json:"match_id" gorm:"type:uuid;uniqueIndex;not null"`
	HomeScore  int       `json:"home_score" gorm:"not null"`
	AwayScore  int       `json:"away_score" gorm:"not null"`
	Version    int       `json:"version" gorm:"not null;default:1"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type MatchResultHistory struct {
	ID            string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MatchResultID string    `json:"match_result_id" gorm:"type:uuid;not null;index"`
	Version       int       `json:"version" gorm:"not null"`
	Changes       string    `json:"changes" gorm:"type:jsonb;not null"`
	CreatedAt     time.Time `json:"created_at"`
}
