package domain

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null;size:255"`
	LogoURL     string         `json:"logo_url" gorm:"size:500"`
	FoundedYear int            `json:"founded_year" gorm:"not null"`
	Address     string         `json:"address" gorm:"not null;type:text"`
	City        string         `json:"city" gorm:"index;not null;size:255"`
	Version     int            `json:"version" gorm:"not null;default:1"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type TeamHistory struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TeamID    string    `json:"team_id" gorm:"type:uuid;not null;index"`
	Version   int       `json:"version" gorm:"not null"`
	Changes   string    `json:"changes" gorm:"type:jsonb;not null"`
	CreatedAt time.Time `json:"created_at"`
}
