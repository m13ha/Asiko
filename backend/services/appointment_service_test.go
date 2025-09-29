package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	repomocks "github.com/m13ha/appointment_master/repository/mocks"
	servicemocks "github.com/m13ha/appointment_master/services/mocks"
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
			expectedError: "Invalid appointment data. Please check your input.",
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
			expectedError: "internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockAppointmentRepo := new(repomocks.AppointmentRepository)
			mockEventNotificationService := new(servicemocks.EventNotificationService)
			tc.setupMock(mockAppointmentRepo, mockEventNotificationService)
			appointmentService := NewAppointmentService(mockAppointmentRepo, mockEventNotificationService)

			// Act
			appointment, err := appointmentService.CreateAppointment(tc.request, userID)

			// Assert
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, appointment)
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
