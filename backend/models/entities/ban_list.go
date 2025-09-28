package entities

import (
	"time"

	"github.com/google/uuid"
)

type BanListEntry struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User        User      `json:"-" gorm:"foreignKey:UserID"`
	BannedEmail string    `json:"banned_email" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
}
