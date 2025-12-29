package notifications

import (
	"context"
	"log"
	"strings"

	"github.com/m13ha/asiko/events"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/repository"
)

// RegisterHandlers subscribes the notification service to relevant events.
func RegisterHandlers(bus events.EventBus, svc NotificationService, bookingRepo repository.BookingRepository) {
	handlers := map[string]func(events.Event) error{
		events.EventBookingCreated: func(event events.Event) error {
			p, ok := event.Data.(events.BookingEventData)
			if !ok || p.Booking == nil {
				return nil
			}
			if strings.ToLower(p.Booking.Status) == entities.BookingStatusPending {
				return nil
			}
			go func() {
				if err := svc.SendBookingConfirmation(p.Booking); err != nil {
					log.Printf("Failed to send booking confirmation: %v", err)
					if bookingRepo != nil {
						bookingRepo.UpdateNotificationStatus(p.Booking.ID, "failed", "email")
					}
				} else if bookingRepo != nil {
					bookingRepo.UpdateNotificationStatus(p.Booking.ID, "sent", "email")
				}
			}()
			return nil
		},
		events.EventBookingCancelled: func(event events.Event) error {
			p, ok := event.Data.(events.BookingEventData)
			if !ok || p.Booking == nil {
				return nil
			}
			go func() {
				if err := svc.SendBookingCancellation(p.Booking); err != nil {
					log.Printf("Failed to send booking cancellation: %v", err)
					if bookingRepo != nil {
						bookingRepo.UpdateNotificationStatus(p.Booking.ID, "failed", "email")
					}
				} else if bookingRepo != nil {
					bookingRepo.UpdateNotificationStatus(p.Booking.ID, "sent", "email")
				}
			}()
			return nil
		},
		events.EventBookingRejected: func(event events.Event) error {
			p, ok := event.Data.(events.BookingEventData)
			if !ok || p.Booking == nil {
				return nil
			}
			go func() {
				if err := svc.SendBookingRejection(p.Booking); err != nil {
					log.Printf("Failed to send booking rejection: %v", err)
					if bookingRepo != nil {
						bookingRepo.UpdateNotificationStatus(p.Booking.ID, "failed", "email")
					}
				} else if bookingRepo != nil {
					bookingRepo.UpdateNotificationStatus(p.Booking.ID, "sent", "email")
				}
			}()
			return nil
		},
		events.EventBookingUpdated: func(event events.Event) error {
			p, ok := event.Data.(events.BookingEventData)
			if !ok || p.Booking == nil {
				return nil
			}
			go func() {
				if err := svc.SendBookingUpdated(p.Booking); err != nil {
					log.Printf("Failed to send booking updated: %v", err)
					if bookingRepo != nil {
						bookingRepo.UpdateNotificationStatus(p.Booking.ID, "failed", "email")
					}
				} else if bookingRepo != nil {
					bookingRepo.UpdateNotificationStatus(p.Booking.ID, "sent", "email")
				}
			}()
			return nil
		},
		events.EventBookingConfirmed: func(event events.Event) error {
			p, ok := event.Data.(events.BookingEventData)
			if !ok || p.Booking == nil {
				return nil
			}
			go func() {
				if err := svc.SendBookingConfirmation(p.Booking); err != nil {
					log.Printf("Failed to send booking confirmation: %v", err)
					if bookingRepo != nil {
						bookingRepo.UpdateNotificationStatus(p.Booking.ID, "failed", "email")
					}
				} else if bookingRepo != nil {
					bookingRepo.UpdateNotificationStatus(p.Booking.ID, "sent", "email")
				}
			}()
			return nil
		},
		events.EventAppointmentCreated: func(event events.Event) error {
			p, ok := event.Data.(events.AppointmentEventData)
			if !ok || p.Appointment == nil {
				return nil
			}
			if strings.TrimSpace(p.RecipientEmail) == "" {
				log.Printf("Skipped appointment created email: missing recipient for %s", p.Appointment.ID)
				return nil
			}
			go func() {
				if err := svc.SendAppointmentCreated(p.Appointment, p.RecipientEmail, p.RecipientName); err != nil {
					log.Printf("Failed to send appointment created: %v", err)
				}
			}()
			return nil
		},
		events.EventAppointmentUpdated: func(event events.Event) error {
			p, ok := event.Data.(events.AppointmentEventData)
			if !ok || p.Appointment == nil {
				return nil
			}
			if strings.TrimSpace(p.RecipientEmail) == "" {
				log.Printf("Skipped appointment updated email: missing recipient for %s", p.Appointment.ID)
				return nil
			}
			go func() {
				if err := svc.SendAppointmentUpdated(p.Appointment, p.RecipientEmail, p.RecipientName); err != nil {
					log.Printf("Failed to send appointment updated: %v", err)
				}
			}()
			return nil
		},
		events.EventAppointmentDeleted: func(event events.Event) error {
			p, ok := event.Data.(events.AppointmentEventData)
			if !ok || p.Appointment == nil {
				return nil
			}
			if strings.TrimSpace(p.RecipientEmail) == "" {
				log.Printf("Skipped appointment deleted email: missing recipient for %s", p.Appointment.ID)
				return nil
			}
			go func() {
				if err := svc.SendAppointmentDeleted(p.Appointment, p.RecipientEmail, p.RecipientName); err != nil {
					log.Printf("Failed to send appointment deleted: %v", err)
				}
			}()
			return nil
		},
	}

	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		handler, ok := handlers[event.Name]
		if !ok {
			return nil
		}
		return handler(event)
	})
}
