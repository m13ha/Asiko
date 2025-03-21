package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AppointmentID uuid.UUID      `json:"appointment_id" gorm:"type:uuid;not null"`
	Appointment   Appointment    `json:"-" gorm:"foreignKey:AppointmentID"`
	AppCode       string         `json:"app_code" gorm:"not null"`
	UserID        *uuid.UUID     `json:"user_id" gorm:"type:uuid"`
	User          User           `json:"-" gorm:"foreignKey:UserID"`
	Name          string         `json:"name" gorm:""`
	Email         string         `json:"email" gorm:""`
	Phone         string         `json:"phone" gorm:""`
	Date          time.Time      `json:"date" gorm:"not null"`
	StartTime     time.Time      `json:"start_time" gorm:"not null"`
	EndTime       time.Time      `json:"end_time" gorm:"not null"`
	Available     bool           `json:"available" gorm:"not null;default:true"`
	AttendeeCount int            `json:"attendee_count" gorm:"default:1"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type BookingRequest struct {
	AppCode       string    `json:"app_code" validate:"required"`
	StartTime     time.Time `json:"start_time" validate:"required"`
	EndTime       time.Time `json:"end_time" validate:"required"`
	Date          time.Time `json:"date" validate:"required"`
	Name          string    `json:"name"`
	Email         string    `json:"email" validate:"omitempty,email"`
	Phone         string    `json:"phone"`
	AttendeeCount int       `json:"attendee_count" validate:"gte=1"`
}

type BookingResponse struct {
	AppCode       string     `json:"app_code"`
	ID            uuid.UUID  `json:"id"`
	AppointmentID uuid.UUID  `json:"appointment_id"`
	UserID        *uuid.UUID `json:"user_id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	Date          time.Time  `json:"date"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       time.Time  `json:"end_time"`
	AttendeeCount int        `json:"attendee_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
