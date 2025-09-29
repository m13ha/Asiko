package entities

import (
	"time"

	"github.com/google/uuid"
)

// Notification represents a notification for a user.

type Notification struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	EventType  string    `json:"event_type" gorm:"not null"`
	Message    string    `json:"message" gorm:"not null"`
	ResourceID uuid.UUID `json:"resource_id" gorm:"type:uuid"`
	IsRead     bool      `json:"is_read" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null;default:now()"`
}
