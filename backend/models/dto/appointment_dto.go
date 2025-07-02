package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/entities"
)

type AppointmentRequest struct {
	Title           string                   `json:"title" validate:"required"`
	StartTime       time.Time                `json:"start_time" validate:"required"`
	EndTime         time.Time                `json:"end_time" validate:"required"`
	BookingDuration int                      `json:"booking_duration" validate:"required,gt=0"`
	StartDate       time.Time                `json:"start_date" validate:"required"`
	EndDate         time.Time                `json:"end_date" validate:"required,gtefield=StartDate"`
	Type            entities.AppointmentType `json:"type" validate:"required,oneof=single group"`
	MaxAttendees    int                      `json:"max_attendees" validate:"gte=1"`
	Description     string                   `json:"description"`
}

type AppointmentResponse struct {
	ID              uuid.UUID                `json:"id"`
	Title           string                   `json:"title"`
	StartTime       time.Time                `json:"start_time"`
	EndTime         time.Time                `json:"end_time"`
	StartDate       time.Time                `json:"start_date"`
	EndDate         time.Time                `json:"end_date"`
	BookingDuration int                      `json:"booking_duration"`
	Type            entities.AppointmentType `json:"type"`
	MaxAttendees    int                      `json:"max_attendees"`
	AppCode         string                   `json:"app_code"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	Description     string                   `json:"description"`
}
