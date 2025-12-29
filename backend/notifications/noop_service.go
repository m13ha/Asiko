package notifications

import (
	"log"

	"github.com/m13ha/asiko/models/entities"
)

// NoopService is a safe fallback when no email provider is configured.
type NoopService struct{}

func NewNoopService() *NoopService {
	return &NoopService{}
}

func (s *NoopService) SendBookingConfirmation(booking *entities.Booking) error {
	log.Printf("notifications: noop booking confirmation for %s", booking.BookingCode)
	return nil
}

func (s *NoopService) SendBookingCancellation(booking *entities.Booking) error {
	log.Printf("notifications: noop booking cancellation for %s", booking.BookingCode)
	return nil
}

func (s *NoopService) SendBookingRejection(booking *entities.Booking) error {
	log.Printf("notifications: noop booking rejection for %s", booking.BookingCode)
	return nil
}

func (s *NoopService) SendBookingUpdated(booking *entities.Booking) error {
	log.Printf("notifications: noop booking updated for %s", booking.BookingCode)
	return nil
}

func (s *NoopService) SendAppointmentCreated(appointment *entities.Appointment, recipientEmail, recipientName string) error {
	log.Printf("notifications: noop appointment created for %s", appointment.Title)
	return nil
}

func (s *NoopService) SendAppointmentUpdated(appointment *entities.Appointment, recipientEmail, recipientName string) error {
	log.Printf("notifications: noop appointment updated for %s", appointment.Title)
	return nil
}

func (s *NoopService) SendAppointmentDeleted(appointment *entities.Appointment, recipientEmail, recipientName string) error {
	log.Printf("notifications: noop appointment deleted for %s", appointment.Title)
	return nil
}

func (s *NoopService) SendVerificationCode(email, code string) error {
	log.Printf("notifications: noop verification email to %s", email)
	return nil
}

func (s *NoopService) SendPasswordResetEmail(email, code string) error {
	log.Printf("notifications: noop password reset email to %s", email)
	return nil
}
