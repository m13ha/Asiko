package services

import (
	"context"
	"log"
	"time"
)

type StatusScheduler interface {
	Start(ctx context.Context)
}

type statusScheduler struct {
	appointmentService AppointmentService
	bookingService     BookingService
	interval           time.Duration
}

func NewStatusScheduler(appointmentService AppointmentService, bookingService BookingService, interval time.Duration) StatusScheduler {
	if interval <= 0 {
		interval = time.Minute
	}
	return &statusScheduler{
		appointmentService: appointmentService,
		bookingService:     bookingService,
		interval:           interval,
	}
}

func (s *statusScheduler) Start(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	go func() {
		s.run(ctx)
	}()
}

func (s *statusScheduler) run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.tick(ctx)

	for {
		select {
		case <-ticker.C:
			s.tick(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *statusScheduler) tick(ctx context.Context) {
	now := time.Now()

	appointmentSummary, err := s.appointmentService.RefreshStatuses(ctx, now)
	if err != nil {
		log.Printf("[StatusScheduler] appointment refresh error: %v", err)
	} else if appointmentSummary.PendingToOngoing+appointmentSummary.Completed > 0 {
		log.Printf(
			"[StatusScheduler] appointment updates — pending→ongoing:%d completed:%d",
			appointmentSummary.PendingToOngoing,
			appointmentSummary.Completed,
		)
	}

	bookingSummary, err := s.bookingService.RefreshBookingStatuses(ctx, now)
	if err != nil {
		log.Printf("[StatusScheduler] booking refresh error: %v", err)
	} else if bookingSummary.Ongoing+bookingSummary.Expired > 0 {
		log.Printf(
			"[StatusScheduler] booking updates — ongoing:%d expired:%d",
			bookingSummary.Ongoing,
			bookingSummary.Expired,
		)
	}
}
