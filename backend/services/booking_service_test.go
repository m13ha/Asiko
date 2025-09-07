package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookAppointment(t *testing.T) {
	userID := uuid.New()
	user := &entities.User{ID: userID, Name: "Test User", Email: "test@example.com"}
	appSlot := &entities.Appointment{ID: uuid.New(), AppCode: "SLOT123", Type: entities.Group, MaxAttendees: 5}
	slot := &entities.Booking{ID: uuid.New(), AppCode: "SLOT123", Available: true, Date: time.Now(), StartTime: time.Now().Add(time.Hour), EndTime: time.Now().Add(90 * time.Minute)}

	t.Run("Success - Book Slot", func(t *testing.T) {
		// Arrange
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockBookingRepo := new(mocks.BookingRepository)
		mockUserRepo := new(mocks.UserRepository)
		bookingService := NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, nil)

		req := requests.BookingRequest{AppCode: "SLOT123", Name: "Test", Email: "test@test.com", Date: slot.Date, StartTime: slot.StartTime, EndTime: slot.EndTime, AttendeeCount: 2}

		mockAppointmentRepo.On("FindAppointmentByAppCode", "SLOT123").Return(appSlot, nil).Once()
		mockBookingRepo.On("FindAvailableSlot", "SLOT123", slot.Date, slot.StartTime).Return(slot, nil).Once()
		mockUserRepo.On("FindByID", userID.String()).Return(user, nil).Once()
		mockBookingRepo.On("Update", mock.AnythingOfType("*entities.Booking")).Return(nil).Once()

		// Act
		booking, err := bookingService.BookAppointment(req, userID.String())

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, booking)
		mockAppointmentRepo.AssertExpectations(t)
		mockBookingRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failure - Appointment not found", func(t *testing.T) {
		// Arrange
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		bookingService := NewBookingService(nil, mockAppointmentRepo, nil, nil)
		validReq := requests.BookingRequest{AppCode: "NOTFOUND", Name: "Guest User", Email: "guest@example.com", Date: time.Now(), StartTime: time.Now(), EndTime: time.Now(), AttendeeCount: 1}
		mockAppointmentRepo.On("FindAppointmentByAppCode", "NOTFOUND").Return(nil, fmt.Errorf("not found")).Once()

		// Act
		_, err := bookingService.BookAppointment(validReq, "")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "Appointment not found.", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
	})

	t.Run("Failure - Slot not available", func(t *testing.T) {
		// Arrange
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockBookingRepo := new(mocks.BookingRepository)
		bookingService := NewBookingService(mockBookingRepo, mockAppointmentRepo, nil, nil)
		validReq := requests.BookingRequest{AppCode: "SLOT123", Name: "Guest User", Email: "guest@example.com", Date: slot.Date, StartTime: slot.StartTime, EndTime: slot.EndTime, AttendeeCount: 1}
		mockAppointmentRepo.On("FindAppointmentByAppCode", "SLOT123").Return(appSlot, nil).Once()
		mockBookingRepo.On("FindAvailableSlot", "SLOT123", slot.Date, slot.StartTime).Return(nil, fmt.Errorf("not found")).Once()

		// Act
		_, err := bookingService.BookAppointment(validReq, "")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "No available slot found.", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
		mockBookingRepo.AssertExpectations(t)
	})
}
