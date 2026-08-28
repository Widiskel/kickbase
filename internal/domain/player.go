package domain

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID           string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TeamID       string         `json:"team_id" gorm:"type:uuid;not null;index"`
	Name         string         `json:"name" gorm:"not null;size:255"`
	Height       float64        `json:"height" gorm:"not null"`
	Weight       float64        `json:"weight" gorm:"not null"`
	Position     string         `json:"position" gorm:"index;not null;size:10"`
	Playstyle    *string        `json:"playstyle,omitempty" gorm:"size:50"`
	JerseyNumber int            `json:"jersey_number" gorm:"index;not null"`
	Version      int            `json:"version" gorm:"not null;default:1"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type PlayerHistory struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PlayerID   string    `json:"player_id" gorm:"type:uuid;not null;index"`
	Version    int       `json:"version" gorm:"not null"`
	Changes    string    `json:"changes" gorm:"type:jsonb;not null"`
	CreatedAt  time.Time `json:"created_at"`
}
