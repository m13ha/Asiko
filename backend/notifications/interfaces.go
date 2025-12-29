package notifications

import "github.com/m13ha/asiko/models/entities"

// NotificationService defines the interface for sending notifications.
type NotificationService interface {
	SendBookingConfirmation(booking *entities.Booking) error
	SendBookingCancellation(booking *entities.Booking) error
	SendBookingRejection(booking *entities.Booking) error
	SendBookingUpdated(booking *entities.Booking) error
	SendAppointmentCreated(appointment *entities.Appointment, recipientEmail, recipientName string) error
	SendAppointmentUpdated(appointment *entities.Appointment, recipientEmail, recipientName string) error
	SendAppointmentDeleted(appointment *entities.Appointment, recipientEmail, recipientName string) error
	SendVerificationCode(email, code string) error
	SendPasswordResetEmail(email, code string) error
}
