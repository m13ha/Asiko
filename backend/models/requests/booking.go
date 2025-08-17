package requests

import (
	"time"
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
