package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/asiko/models/entities"
)

// Event name constants
const (
	EventBookingCreated   = "booking.created"
	EventBookingCancelled = "booking.cancelled"
	EventBookingUpdated   = "booking.updated"
	EventBookingRejected  = "booking.rejected"
	EventBookingConfirmed = "booking.confirmed"
	EventAppointmentCreated = "appointment.created"
	EventAppointmentUpdated = "appointment.updated"
	EventAppointmentDeleted = "appointment.deleted"
)

// Event is the standard envelope for all events published on the bus.
type Event struct {
	Name string
	Data interface{}
}

// Option 1: Pass the full entity.
// Option 2: Define specific payloads.
// For now, passing the full entity is convenient for avoiding mapping overhead,
// as the consumers (NotificationService) operate on entities.
// However, to keep 'events' package decoupled from 'models' (if we wanted strict layering), we'd use DTOs.
// Given the current monolith structure, importing entities is acceptable and pragmatic.

type BookingEventData struct {
	Booking          *entities.Booking
	OwnerID          uuid.UUID
	AppointmentTitle string
	RecipientEmail   string
	RecipientName    string
}

type AppointmentEventData struct {
	Appointment      *entities.Appointment
	OwnerID          uuid.UUID
	AppointmentTitle string
	RecipientEmail   string
	RecipientName    string
}

// If we need stricter decoupling later, we can use these DTOs:
type BookingDTO struct {
	ID            uuid.UUID
	BookingCode   string
	AppCode       string
	Name          string
	Email         string
	AppointmentID uuid.UUID
	OwnerID       uuid.UUID // Appointment Owner
	StartDate     time.Time
}
