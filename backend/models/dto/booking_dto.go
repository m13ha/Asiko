package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookingRequest struct {
	AppCode       string    `json:"app_code" validate:"required"`
	StartTime     time.Time `json:"start_time" validate:"required"`
	EndTime       time.Time `json:"end_time" validate:"required"`
	Date          time.Time `json:"date" validate:"required"`
	Name          string    `json:"name"`
	Email         string    `json:"email" validate:"omitempty,email"`
	Phone         string    `json:"phone"`
	AttendeeCount int       `json:"attendee_count" validate:"gte=1"`
	Description   string    `json:"description"`
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
	Description   string     `json:"description"`
}
