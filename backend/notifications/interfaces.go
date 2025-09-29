package notifications

import "github.com/m13ha/appointment_master/models/entities"

// NotificationService defines the interface for sending notifications.
type NotificationService interface {
	SendBookingConfirmation(booking *entities.Booking) error
	SendBookingCancellation(booking *entities.Booking) error
	SendBookingRejection(booking *entities.Booking) error
	SendVerificationCode(email, code string) error
}
