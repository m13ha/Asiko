package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/m13ha/asiko/services"
	servicesmocks "github.com/m13ha/asiko/services/mocks"
	"github.com/stretchr/testify/mock"
)

func TestStatusSchedulerTickInvokesServices(t *testing.T) {
	mockAppointmentService := new(servicesmocks.AppointmentService)
	mockBookingService := new(servicesmocks.BookingService)

	callDone := make(chan struct{}, 2)
	mockAppointmentService.
		On("RefreshStatuses", mock.Anything, mock.AnythingOfType("time.Time")).
		Run(func(args mock.Arguments) { callDone <- struct{}{} }).
		Return(services.StatusRefreshSummary{PendingToOngoing: 1, Completed: 2}, nil).
		Once()
	mockBookingService.
		On("RefreshBookingStatuses", mock.Anything, mock.AnythingOfType("time.Time")).
		Run(func(args mock.Arguments) { callDone <- struct{}{} }).
		Return(services.BookingStatusRefreshSummary{Ongoing: 3, Expired: 4}, nil).
		Once()

	scheduler := services.NewStatusScheduler(mockAppointmentService, mockBookingService, time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	scheduler.Start(ctx)

	timeout := time.After(500 * time.Millisecond)
	for i := 0; i < 2; i++ {
		select {
		case <-callDone:
		case <-timeout:
			t.Fatal("scheduler did not invoke services in time")
		}
	}

	mockAppointmentService.AssertExpectations(t)
	mockBookingService.AssertExpectations(t)
}
