package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:100"`
	Password  string         `json:"-" gorm:"not null;size:255"`
	Name      string         `json:"name" gorm:"not null;size:255"`
	Role      string         `json:"role" gorm:"not null;default:'admin';size:50"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
