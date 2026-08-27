package domain

import (
	"time"
)

type Goal struct {
	ID             string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MatchResultID  string    `json:"match_result_id" gorm:"type:uuid;not null;index"`
	PlayerID       string    `json:"player_id" gorm:"type:uuid;not null;index"`
	GoalTime       string    `json:"goal_time" gorm:"not null;size:10"`
	CreatedAt      time.Time `json:"created_at"`
}

type GoalHistory struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GoalID    string    `json:"goal_id" gorm:"type:uuid;not null;index"`
	Version   int       `json:"version" gorm:"not null"`
	Changes   string    `json:"changes" gorm:"type:jsonb;not null"`
	CreatedAt time.Time `json:"created_at"`
}
