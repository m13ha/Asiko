package services

import (
	"context"
	"fmt"
	"log"

	"github.com/m13ha/asiko/events"
)

// RegisterInternalHandlers subscribes the internal notification service to relevant events.
func RegisterInternalHandlers(bus events.EventBus, svc EventNotificationService) {
	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		if event.Name != events.EventBookingCreated {
			return nil
		}
		p, ok := event.Data.(events.BookingEventData)
		if !ok || p.Booking == nil {
			return nil
		}

		message := fmt.Sprintf("New booking by %s for your appointment %s.", p.Booking.Name, p.AppointmentTitle)
		// Internal Notification needs to go to the OWNER of the appointment
		if err := svc.CreateEventNotification(p.OwnerID, "BOOKING_CREATED", message, p.Booking.ID); err != nil {
			log.Printf("Failed to create internal notification: %v", err)
		}
		return nil
	})

	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		if event.Name != events.EventAppointmentCreated {
			return nil
		}
		p, ok := event.Data.(events.AppointmentEventData)
		if !ok || p.Appointment == nil {
			return nil
		}

		message := fmt.Sprintf("New appointment '%s' created.", p.AppointmentTitle)
		if err := svc.CreateEventNotification(p.OwnerID, "APPOINTMENT_CREATED", message, p.Appointment.ID); err != nil {
			log.Printf("Failed to create internal notification: %v", err)
		}
		return nil
	})

	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		if event.Name != events.EventAppointmentUpdated {
			return nil
		}
		p, ok := event.Data.(events.AppointmentEventData)
		if !ok || p.Appointment == nil {
			return nil
		}

		message := fmt.Sprintf("Appointment '%s' was updated.", p.AppointmentTitle)
		if err := svc.CreateEventNotification(p.OwnerID, "APPOINTMENT_UPDATED", message, p.Appointment.ID); err != nil {
			log.Printf("Failed to create internal notification: %v", err)
		}
		return nil
	})

	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		if event.Name != events.EventAppointmentDeleted {
			return nil
		}
		p, ok := event.Data.(events.AppointmentEventData)
		if !ok || p.Appointment == nil {
			return nil
		}

		message := fmt.Sprintf("Appointment '%s' was deleted.", p.AppointmentTitle)
		if err := svc.CreateEventNotification(p.OwnerID, "APPOINTMENT_DELETED", message, p.Appointment.ID); err != nil {
			log.Printf("Failed to create internal notification: %v", err)
		}
		return nil
	})

	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		if event.Name != events.EventBookingCancelled {
			return nil
		}
		p, ok := event.Data.(events.BookingEventData)
		if !ok || p.Booking == nil {
			return nil
		}

		message := fmt.Sprintf("Booking %s was cancelled.", p.Booking.BookingCode)
		if err := svc.CreateEventNotification(p.OwnerID, "BOOKING_CANCELLED", message, p.Booking.ID); err != nil {
			log.Printf("Failed to create internal notification: %v", err)
		}
		return nil
	})

	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		if event.Name != events.EventBookingUpdated {
			return nil
		}
		p, ok := event.Data.(events.BookingEventData)
		if !ok || p.Booking == nil {
			return nil
		}

		message := fmt.Sprintf("Booking %s was updated.", p.Booking.BookingCode)
		if err := svc.CreateEventNotification(p.OwnerID, "BOOKING_UPDATED", message, p.Booking.ID); err != nil {
			log.Printf("Failed to create internal notification: %v", err)
		}
		return nil
	})

	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		if event.Name != events.EventBookingRejected {
			return nil
		}
		p, ok := event.Data.(events.BookingEventData)
		if !ok || p.Booking == nil {
			return nil
		}
		if p.Booking.UserID == nil {
			return nil
		}

		message := fmt.Sprintf("Your booking %s for %s was rejected.", p.Booking.BookingCode, p.AppointmentTitle)
		if err := svc.CreateEventNotification(*p.Booking.UserID, "BOOKING_REJECTED", message, p.Booking.ID); err != nil {
			log.Printf("Failed to create attendee notification: %v", err)
		}
		return nil
	})
}
