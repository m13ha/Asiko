package responses

import (
	"time"

	"github.com/google/uuid"
)

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
