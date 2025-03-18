package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AppointmentID uuid.UUID      `json:"appointment_id" gorm:"type:uuid;not null"`
	Appointment   Appointment    `json:"appointment" gorm:"foreignKey:AppointmentID"`
	UserID        *uuid.UUID     `json:"user_id" gorm:"type:uuid"` // Nullable for guest bookings
	User          User           `json:"-" gorm:"foreignKey:UserID"`
	GuestName     string         `json:"guest_name" gorm:""`
	GuestEmail    string         `json:"guest_email" gorm:""`
	GuestPhone    string         `json:"guest_phone" gorm:""`
	Date          time.Time      `json:"date" gorm:"not null"`
	StartTime     time.Time      `json:"start_time" gorm:"not null"`
	EndTime       time.Time      `json:"end_time" gorm:"not null"`
	Available     bool           `json:"available" gorm:"not null;default:true"`
	AttendeeCount int            `json:"attendee_count" gorm:"default:1"` // For group bookings
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type BookingRequest struct {
	AppointmentID uuid.UUID `json:"appointment_id" validate:"required"`
	StartTime     time.Time `json:"start_time" validate:"required"`
	EndTime       time.Time `json:"end_time" validate:"required"`
	Date          time.Time `json:"date" validate:"required"`
	// For guest bookings
	GuestName     string    `json:"guest_name"`
	GuestEmail    string    `json:"guest_email" validate:"omitempty,email"`
	GuestPhone    string    `json:"guest_phone"`
	AttendeeCount int       `json:"attendee_count" validate:"gte=1"`
	UserID        uuid.UUID `json:"user_id"`
}

type BookingResponse struct {
	ID            uuid.UUID  `json:"id"`
	AppointmentID uuid.UUID  `json:"appointment_id"`
	UserID        *uuid.UUID `json:"user_id"`
	GuestName     string     `json:"guest_name"`
	GuestEmail    string     `json:"guest_email"`
	GuestPhone    string     `json:"guest_phone"`
	Date          time.Time  `json:"date"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       time.Time  `json:"end_time"`
	AttendeeCount int        `json:"attendee_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
