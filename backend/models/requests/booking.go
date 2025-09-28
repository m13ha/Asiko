package requests

import (
	"time"

	myerrors "github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/utils"
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
	DeviceToken   string    `json:"device_token,omitempty"`
}

func (req *BookingRequest) Validate() error {
	if err := utils.Validate(req); err != nil {
		return myerrors.NewUserError("Invalid booking data. Please check your input.")
	}

	if req.Name == "" || (req.Email == "" && req.Phone == "") {
		return myerrors.NewUserError("Name and either email or phone are required for guest bookings.")
	}
	return nil
}
