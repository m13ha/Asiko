package responses

import (
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/asiko/models/entities"
)

type AppointmentResponse struct {
	ID              uuid.UUID                  `json:"id"`
	Title           string                     `json:"title"`
	StartTime       time.Time                  `json:"start_time"`
	EndTime         time.Time                  `json:"end_time"`
	StartDate       time.Time                  `json:"start_date"`
	EndDate         time.Time                  `json:"end_date"`
	BookingDuration int                        `json:"booking_duration"`
	Type            entities.AppointmentType   `json:"type"`
	MaxAttendees    int                        `json:"max_attendees"`
	AppCode         string                     `json:"app_code"`
	Status          entities.AppointmentStatus `json:"status"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
	Description     string                     `json:"description"`
}
