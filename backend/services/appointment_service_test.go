package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/asiko/events"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	repomocks "github.com/m13ha/asiko/repository/mocks"
	services "github.com/m13ha/asiko/services"
	servicemocks "github.com/m13ha/asiko/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event events.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(handler events.EventHandler) {
	m.Called(handler)
}

func TestCreateAppointment(t *testing.T) {
	userID := uuid.New()

	testCases := []struct {
		name          string
		request       requests.AppointmentRequest
		setupMock     func(mockRepo *repomocks.AppointmentRepository, mockEventBus *MockEventBus, mockEventNotificationService *servicemocks.EventNotificationService)
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
			setupMock: func(mockRepo *repomocks.AppointmentRepository, mockEventBus *MockEventBus, mockEventNotificationService *servicemocks.EventNotificationService) {
				mockRepo.On("Create", mock.AnythingOfType("*entities.Appointment")).Return(nil).Once()
				mockEventBus.On("Publish", mock.Anything, mock.MatchedBy(func(event events.Event) bool {
					return event.Name == events.EventAppointmentCreated
				})).Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name: "Failure - Validation Error",
			request: requests.AppointmentRequest{
				Title: "", // Invalid title
			},
			setupMock: func(mockRepo *repomocks.AppointmentRepository, mockEventBus *MockEventBus, mockEventNotificationService *servicemocks.EventNotificationService) {
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
			setupMock: func(mockRepo *repomocks.AppointmentRepository, mockEventBus *MockEventBus, mockEventNotificationService *servicemocks.EventNotificationService) {
				mockRepo.On("Create", mock.AnythingOfType("*entities.Appointment")).Return(fmt.Errorf("db error")).Once()
			},
			expectedError: "INTERNAL_ERROR: db error (caused by: db error)",
		},
		{
			name: "Success - Overnight Group Appointment",
			request: func() requests.AppointmentRequest {
				startDate := time.Date(2025, time.January, 10, 0, 0, 0, 0, time.UTC)
				endDate := startDate.AddDate(0, 0, 1)
				startTime := time.Date(2025, time.January, 10, 22, 0, 0, 0, time.UTC)
				endTime := time.Date(2025, time.January, 11, 2, 0, 0, 0, time.UTC)
				return requests.AppointmentRequest{
					Title:           "Overnight Group",
					StartTime:       startTime,
					EndTime:         endTime,
					StartDate:       startDate,
					EndDate:         endDate,
					BookingDuration: 30,
					Type:            entities.Group,
					MaxAttendees:    5,
				}
			}(),
			setupMock: func(mockRepo *repomocks.AppointmentRepository, mockEventBus *MockEventBus, mockEventNotificationService *servicemocks.EventNotificationService) {
				mockRepo.On("Create", mock.AnythingOfType("*entities.Appointment")).Return(nil).Once()
				mockEventBus.On("Publish", mock.Anything, mock.MatchedBy(func(event events.Event) bool {
					return event.Name == events.EventAppointmentCreated
				})).Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name: "Failure - Party More Than One Day",
			request: func() requests.AppointmentRequest {
				startDate := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
				endDate := startDate.AddDate(0, 0, 2)
				startTime := time.Date(2025, time.February, 1, 20, 0, 0, 0, time.UTC)
				endTime := time.Date(2025, time.February, 3, 2, 0, 0, 0, time.UTC)
				return requests.AppointmentRequest{
					Title:           "Long Party",
					StartTime:       startTime,
					EndTime:         endTime,
					StartDate:       startDate,
					EndDate:         endDate,
					BookingDuration: 60,
					Type:            entities.Party,
					MaxAttendees:    50,
				}
			}(),
			setupMock: func(mockRepo *repomocks.AppointmentRepository, mockEventBus *MockEventBus, mockEventNotificationService *servicemocks.EventNotificationService) {},
			expectedError: "VALIDATION_FAILED: Party appointments cannot span more than one day.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockAppointmentRepo := new(repomocks.AppointmentRepository)
			mockBookingRepo := new(repomocks.BookingRepository)
			mockUserRepo := new(repomocks.UserRepository)
			mockEventBus := new(MockEventBus)
			mockEventNotificationService := new(servicemocks.EventNotificationService)
			mockUserRepo.On("FindByID", userID.String()).Return(&entities.User{ID: userID, Name: "Owner", Email: "owner@example.com"}, nil).Maybe()
			tc.setupMock(mockAppointmentRepo, mockEventBus, mockEventNotificationService)
			appointmentService := services.NewAppointmentService(mockAppointmentRepo, mockBookingRepo, mockUserRepo, mockEventBus, mockEventNotificationService, nil)

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
			mockEventBus.AssertExpectations(t)
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

	mockEventBus := new(MockEventBus)
	svc := services.NewAppointmentService(mockAppointmentRepo, new(repomocks.BookingRepository), new(repomocks.UserRepository), mockEventBus, mockEventNotificationService, nil)

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

	mockEventBus := new(MockEventBus)
	svc := services.NewAppointmentService(mockAppointmentRepo, new(repomocks.BookingRepository), new(repomocks.UserRepository), mockEventBus, mockEventNotificationService, nil)

	summary, err := svc.RefreshStatuses(ctx, now)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), summary.PendingToOngoing)
	assert.Equal(t, int64(1), summary.Completed)
	mockAppointmentRepo.AssertExpectations(t)
}
