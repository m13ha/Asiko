package services

import (
	"context"
	"log"
	"time"
)

type BookingStatusScheduler interface {
	Start(ctx context.Context)
}

type bookingStatusScheduler struct {
	service  BookingService
	interval time.Duration
}

func NewBookingStatusScheduler(service BookingService, interval time.Duration) BookingStatusScheduler {
	if interval <= 0 {
		interval = time.Minute
	}
	return &bookingStatusScheduler{service: service, interval: interval}
}

func (s *bookingStatusScheduler) Start(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	go func() {
		s.run(ctx)
	}()
}

func (s *bookingStatusScheduler) run(ctx context.Context) {
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

func (s *bookingStatusScheduler) tick(ctx context.Context) {
	summary, err := s.service.RefreshBookingStatuses(ctx, time.Now())
	if err != nil {
		log.Printf("[BookingStatusScheduler] refresh error: %v", err)
		return
	}

	if summary.Ongoing+summary.Expired > 0 {
		log.Printf("[BookingStatusScheduler] status updates — ongoing:%d expired:%d", summary.Ongoing, summary.Expired)
	}
}
