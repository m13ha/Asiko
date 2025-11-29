package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	repomocks "github.com/m13ha/asiko/repository/mocks"
	services "github.com/m13ha/asiko/services"
	servicemocks "github.com/m13ha/asiko/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAppointment(t *testing.T) {
	userID := uuid.New()

	testCases := []struct {
		name          string
		request       requests.AppointmentRequest
		setupMock     func(mockRepo *repomocks.AppointmentRepository, mockEventNotificationService *servicemocks.EventNotificationService)
		expectedError string
	}{
		{
			name: "Success",
			request: requests.AppointmentRequest{
				Title:           "Test Appointment",
				StartTime:       time.Now().Add(time.Hour * 24),
				EndTime:         time.Now().Add(time.Hour * 25),
				StartDate:       time.Now().Add(time.Hour * 24),
				EndDate:         time.Now().Add(time.Hour * 48),
				BookingDuration: 30,
				Type:            entities.Single,
				MaxAttendees:    1,
			},
			setupMock: func(mockRepo *repomocks.AppointmentRepository, mockEventNotificationService *servicemocks.EventNotificationService) {
				mockRepo.On("Create", mock.AnythingOfType("*entities.Appointment")).Return(nil).Once()
				mockEventNotificationService.On("CreateEventNotification", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name: "Failure - Validation Error",
			request: requests.AppointmentRequest{
				Title: "", // Invalid title
			},
			setupMock: func(mockRepo *repomocks.AppointmentRepository, mockEventNotificationService *servicemocks.EventNotificationService) {
			},
			expectedError: "USER_ERROR: Invalid appointment data. Please check your input.",
		},
		{
			name: "Failure - Repository Error",
			request: requests.AppointmentRequest{
				Title:           "Test Appointment",
				StartTime:       time.Now().Add(time.Hour * 24),
				EndTime:         time.Now().Add(time.Hour * 25),
				StartDate:       time.Now().Add(time.Hour * 24),
				EndDate:         time.Now().Add(time.Hour * 48),
				BookingDuration: 30,
				Type:            entities.Single,
				MaxAttendees:    1,
			},
			setupMock: func(mockRepo *repomocks.AppointmentRepository, mockEventNotificationService *servicemocks.EventNotificationService) {
				mockRepo.On("Create", mock.AnythingOfType("*entities.Appointment")).Return(fmt.Errorf("db error")).Once()
			},
			expectedError: "INTERNAL_ERROR: db error (caused by: db error)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockAppointmentRepo := new(repomocks.AppointmentRepository)
			mockEventNotificationService := new(servicemocks.EventNotificationService)
			tc.setupMock(mockAppointmentRepo, mockEventNotificationService)
			appointmentService := services.NewAppointmentService(mockAppointmentRepo, mockEventNotificationService)

			// Act
			appointment, err := appointmentService.CreateAppointment(tc.request, userID)

			// Assert
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, appointment)
				assert.Equal(t, entities.AppointmentStatusPending, appointment.Status)
			} else {
				assert.Error(t, err)
				assert.Nil(t, appointment)
				assert.Equal(t, tc.expectedError, err.Error())
			}
			mockAppointmentRepo.AssertExpectations(t)
			mockEventNotificationService.AssertExpectations(t)
		})
	}
}

func TestCancelAppointment(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()
	appointmentID := uuid.New()
	appointment := &entities.Appointment{ID: appointmentID, OwnerID: ownerID, Title: "Demo", Status: entities.AppointmentStatusPending}

	mockAppointmentRepo := new(repomocks.AppointmentRepository)
	mockEventNotificationService := new(servicemocks.EventNotificationService)
	mockAppointmentRepo.On("FindByIDAndOwner", ctx, appointmentID, ownerID).Return(appointment, nil).Once()
	mockAppointmentRepo.On("UpdateStatus", ctx, appointmentID, entities.AppointmentStatusCanceled).Return(nil).Once()
	mockEventNotificationService.On("CreateEventNotification", ownerID, "APPOINTMENT_CANCELED", mock.AnythingOfType("string"), appointmentID).Return(nil).Once()

	svc := services.NewAppointmentService(mockAppointmentRepo, mockEventNotificationService)

	result, err := svc.CancelAppointment(ctx, appointmentID, ownerID)

	assert.NoError(t, err)
	assert.Equal(t, entities.AppointmentStatusCanceled, result.Status)
	mockAppointmentRepo.AssertExpectations(t)
	mockEventNotificationService.AssertExpectations(t)
}

func TestRefreshStatuses(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	mockAppointmentRepo := new(repomocks.AppointmentRepository)
	mockEventNotificationService := new(servicemocks.EventNotificationService)
	mockAppointmentRepo.On("MarkAppointmentsOngoing", ctx, now).Return(int64(2), nil).Once()
	mockAppointmentRepo.On("MarkAppointmentsCompleted", ctx, now).Return(int64(1), nil).Once()
	mockAppointmentRepo.On("MarkAppointmentsExpired", ctx, now).Return(int64(3), nil).Once()

	svc := services.NewAppointmentService(mockAppointmentRepo, mockEventNotificationService)

	summary, err := svc.RefreshStatuses(ctx, now)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), summary.PendingToOngoing)
	assert.Equal(t, int64(1), summary.Completed)
	assert.Equal(t, int64(3), summary.Expired)
	mockAppointmentRepo.AssertExpectations(t)
}
