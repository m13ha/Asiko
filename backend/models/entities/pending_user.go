package entities

import (
	"time"

	"github.com/google/uuid"
)

// PendingUser represents a user who has registered but not yet verified their email.

type PendingUser struct {
	ID                        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name                      string    `json:"name" gorm:"not null"`
	Email                     string    `json:"email" gorm:"not null;unique"`
	HashedPassword            string    `json:"-" gorm:"not null"`
	PhoneNumber               *string   `json:"phone_number"`
	VerificationCode          string    `json:"-" gorm:"not null"`
	VerificationCodeExpiresAt time.Time `json:"-" gorm:"not null"`
	CreatedAt                 time.Time `json:"created_at" gorm:"not null;default:now()"`
}
