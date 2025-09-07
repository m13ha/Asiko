package requests

import (
	"time"

	myerrors "github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/utils"
)

type AppointmentRequest struct {
	Title           string                   `json:"title" validate:"required"`
	StartTime       time.Time                `json:"start_time" validate:"required"`
	EndTime         time.Time                `json:"end_time" validate:"required"`
	BookingDuration int                      `json:"booking_duration" validate:"required,gt=0"`
	StartDate       time.Time                `json:"start_date" validate:"required"`
	EndDate         time.Time                `json:"end_date" validate:"required,gtefield=StartDate"`
	Type            entities.AppointmentType `json:"type" validate:"required,oneof=single group party"`
	MaxAttendees    int                      `json:"max_attendees" validate:"gte=1"`
	Description     string                   `json:"description"`
}

func (req *AppointmentRequest) Validate() error {

	if err := utils.Validate(req); err != nil {
		return myerrors.NewUserError("Invalid appointment data. Please check your input.")
	}

	if req.EndTime.Before(req.StartTime) {
		return myerrors.NewUserError("End time cannot be before start time.")
	}

	if req.EndDate.Before(req.StartDate) {
		return myerrors.NewUserError("End date cannot be before start date.")
	}

	duration := req.EndTime.Sub(req.StartTime)
	if duration.Minutes() < float64(req.BookingDuration) {
		return myerrors.NewUserError("Booking duration exceeds available time window.")
	}

	return nil
}
