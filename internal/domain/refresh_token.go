package domain

import (
	"time"
)

type RefreshToken struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null;size:500"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Revoked   bool      `json:"revoked" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
