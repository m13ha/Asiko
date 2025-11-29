package requests

import (
	"time"

	serviceerrors "github.com/m13ha/asiko/errors/serviceerrors"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/utils"
)

type AppointmentRequest struct {
	Title             string                     `json:"title" validate:"required"`
	StartTime         time.Time                  `json:"start_time" validate:"required"`
	EndTime           time.Time                  `json:"end_time" validate:"required"`
	BookingDuration   int                        `json:"booking_duration" validate:"required,gt=0"`
	StartDate         time.Time                  `json:"start_date" validate:"required"`
	EndDate           time.Time                  `json:"end_date" validate:"required,gtefield=StartDate"`
	Type              entities.AppointmentType   `json:"type" validate:"required,oneof=single group party"`
	MaxAttendees      int                        `json:"max_attendees" validate:"gte=1"`
	Description       string                     `json:"description"`
	AntiScalpingLevel entities.AntiScalpingLevel `json:"anti_scalping_level,omitempty" validate:"omitempty,oneof=none standard strict"`
}

func (req *AppointmentRequest) Validate() error {

	if err := utils.Validate(req); err != nil {
		return serviceerrors.UserError("Invalid appointment data. Please check your input.")
	}

	startClock := normalizeClock(req.StartTime)
	endClock := normalizeClock(req.EndTime)

	if !endClock.After(startClock) {
		return serviceerrors.ValidationError("End time cannot be before start time.")
	}

	if req.EndDate.Before(req.StartDate) {
		return serviceerrors.ValidationError("End date cannot be before start date.")
	}

	duration := endClock.Sub(startClock)
	if duration.Minutes() < float64(req.BookingDuration) {
		return serviceerrors.ValidationError("Booking duration exceeds available time window.")
	}

	// Align times to the start date for downstream processing
	req.StartTime = time.Date(
		req.StartDate.Year(), req.StartDate.Month(), req.StartDate.Day(),
		req.StartTime.Hour(), req.StartTime.Minute(), req.StartTime.Second(), req.StartTime.Nanosecond(),
		req.StartTime.Location(),
	)
	req.EndTime = time.Date(
		req.StartDate.Year(), req.StartDate.Month(), req.StartDate.Day(),
		req.EndTime.Hour(), req.EndTime.Minute(), req.EndTime.Second(), req.EndTime.Nanosecond(),
		req.EndTime.Location(),
	)

	return nil
}

func normalizeClock(t time.Time) time.Time {
	return time.Date(2000, time.January, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}
