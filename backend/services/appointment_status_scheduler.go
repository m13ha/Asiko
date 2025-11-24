package services

import (
	"context"
	"log"
	"time"
)

type AppointmentStatusScheduler interface {
	Start(ctx context.Context)
}

type appointmentStatusScheduler struct {
	service  AppointmentService
	interval time.Duration
}

func NewAppointmentStatusScheduler(service AppointmentService, interval time.Duration) AppointmentStatusScheduler {
	if interval <= 0 {
		interval = time.Minute
	}
	return &appointmentStatusScheduler{service: service, interval: interval}
}

func (s *appointmentStatusScheduler) Start(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	go func() {
		s.run(ctx)
	}()
}

func (s *appointmentStatusScheduler) run(ctx context.Context) {
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

func (s *appointmentStatusScheduler) tick(ctx context.Context) {
	summary, err := s.service.RefreshStatuses(ctx, time.Now())
	if err != nil {
		log.Printf("[AppointmentStatusScheduler] refresh error: %v", err)
		return
	}

	if summary.PendingToOngoing+summary.Completed+summary.Expired > 0 {
		log.Printf("[AppointmentStatusScheduler] status updates — pending→ongoing:%d completed:%d expired:%d", summary.PendingToOngoing, summary.Completed, summary.Expired)
	}
}
