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
	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		if event.Name != events.EventBookingCreated {
			return nil
		}
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
	})

	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		if event.Name != events.EventBookingCancelled {
			return nil
		}
		p, ok := event.Data.(events.BookingEventData)
		if !ok || p.Booking == nil {
			return nil
		}
		go func() {
			if err := svc.SendBookingCancellation(p.Booking); err != nil {
				log.Printf("Failed to send booking cancellation: %v", err)
			}
		}()
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
	})

	bus.Subscribe(func(ctx context.Context, event events.Event) error {
		if event.Name != events.EventBookingConfirmed {
			return nil
		}
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
	})
}
