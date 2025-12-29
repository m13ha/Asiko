package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	repoerrors "github.com/m13ha/asiko/errors/repoerrors"
	"github.com/m13ha/asiko/events"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	repomocks "github.com/m13ha/asiko/repository/mocks"
	services "github.com/m13ha/asiko/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBookAppointment(t *testing.T) {
	userID := uuid.New()
	user := &entities.User{ID: userID, Name: "Test User", Email: "test@example.com"}
	appSlot := &entities.Appointment{ID: uuid.New(), AppCode: "SLOT123", Type: entities.Group, MaxAttendees: 5, AntiScalpingLevel: entities.ScalpingNone}
	slotDate := time.Now()
	slotStart := slotDate.Add(time.Hour)
	slotEnd := slotDate.Add(90 * time.Minute)
	newSlot := func() *entities.Booking {
		return &entities.Booking{
			ID:          uuid.New(),
			AppCode:     "SLOT123",
			Available:   true,
			IsSlot:      true,
			Capacity:    5,
			SeatsBooked: 0,
			Date:        slotDate,
			StartTime:   slotStart,
			EndTime:     slotEnd,
		}
	}

	t.Run("Success - Book Slot", func(t *testing.T) {
		// Arrange
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		stubAppointmentWithTx(mockAppointmentRepo)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockEventBus := new(MockEventBus)
		stubBanListNotFound(mockBanListRepo)

		db, sqlMock, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := services.NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, gormDB)

		slot := newSlot()
		req := requests.BookingRequest{AppCode: "SLOT123", Name: "Test", Email: "test@test.com", Date: slot.Date, StartTime: slot.StartTime, EndTime: slot.EndTime, AttendeeCount: 2}

		mockAppointmentRepo.On("FindAppointmentByAppCode", "SLOT123").Return(appSlot, nil).Once()
		mockAppointmentRepo.On("FindAndLock", "SLOT123", mock.AnythingOfType("*gorm.DB")).Return(appSlot, nil).Once()
		mockUserRepo.On("FindByID", userID.String()).Return(user, nil).Once()
		// Transaction expectations
		sqlMock.ExpectBegin()
		// Repository uses WithTx + FindAndLockAvailableSlot in transaction
		mockBookingRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookingRepo).Once()
		mockBookingRepo.On("FindAndLockAvailableSlot", "SLOT123", slot.Date, slot.StartTime).Return(slot, nil).Once()
		mockBookingRepo.On("Create", mock.AnythingOfType("*entities.Booking")).Return(nil).Once()
		mockBookingRepo.On("Update", mock.AnythingOfType("*entities.Booking")).Return(nil).Once()
		sqlMock.ExpectCommit()

		// Event Bus Expectations
		mockEventBus.On("Publish", mock.Anything, mock.MatchedBy(func(event events.Event) bool {
			return event.Name == events.EventBookingCreated
		})).Return(nil).Once()

		// Act
		booking, err := bookingService.BookAppointment(req, userID.String())

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, booking)
		assert.False(t, booking.IsSlot)
		assert.Equal(t, 2, booking.AttendeeCount)
		assert.Equal(t, 2, slot.SeatsBooked)
		assert.Equal(t, 3, slot.AttendeeCount)
		assert.True(t, slot.Available)
		mockAppointmentRepo.AssertExpectations(t)
		mockBookingRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockEventBus.AssertExpectations(t)
	})

	t.Run("Failure - Group Capacity Exceeded", func(t *testing.T) {
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		stubAppointmentWithTx(mockAppointmentRepo)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockEventBus := new(MockEventBus)
		stubBanListNotFound(mockBanListRepo)
		db, sqlMock, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := services.NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, gormDB)

		req := requests.BookingRequest{AppCode: "SLOT123", Name: "Test", Email: "test@test.com", Date: slotDate, StartTime: slotStart, EndTime: slotEnd, AttendeeCount: 4}
		mockAppointmentRepo.On("FindAppointmentByAppCode", "SLOT123").Return(appSlot, nil).Once()
		mockAppointmentRepo.On("FindAndLock", "SLOT123", mock.AnythingOfType("*gorm.DB")).Return(appSlot, nil).Once()
		mockUserRepo.On("FindByID", userID.String()).Return(user, nil).Once()
		sqlMock.ExpectBegin()
		mockBookingRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookingRepo).Once()
		partialSlot := *newSlot()
		partialSlot.SeatsBooked = 3
		partialSlot.AttendeeCount = 2
		mockBookingRepo.On("FindAndLockAvailableSlot", "SLOT123", partialSlot.Date, partialSlot.StartTime).Return(&partialSlot, nil).Once()
		sqlMock.ExpectRollback()

		_, err := bookingService.BookAppointment(req, userID.String())

		assert.Error(t, err)
		assert.Equal(t, "BOOKING_CAPACITY_EXCEEDED: not enough capacity for this slot", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
		mockBookingRepo.AssertExpectations(t)
	})

	t.Run("Failure - Appointment not found", func(t *testing.T) {
		// Arrange
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockEventBus := new(MockEventBus)
		stubBanListNotFound(mockBanListRepo)
		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := services.NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, gormDB)
		validReq := requests.BookingRequest{AppCode: "NOTFOUND", Name: "Guest User", Email: "guest@example.com", Date: time.Now(), StartTime: time.Now(), EndTime: time.Now(), AttendeeCount: 1}
		mockAppointmentRepo.On("FindAppointmentByAppCode", "NOTFOUND").Return(nil, fmt.Errorf("not found")).Once()

		// Act
		_, err := bookingService.BookAppointment(validReq, "")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
	})

	t.Run("Failure - Slot not available", func(t *testing.T) {
		// Arrange
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		stubAppointmentWithTx(mockAppointmentRepo)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockEventBus := new(MockEventBus)
		stubBanListNotFound(mockBanListRepo)
		db, sqlMock, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := services.NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, gormDB)
		slot := newSlot()
		validReq := requests.BookingRequest{AppCode: "SLOT123", Name: "Guest User", Email: "guest@example.com", Date: slot.Date, StartTime: slot.StartTime, EndTime: slot.EndTime, AttendeeCount: 1}
		mockAppointmentRepo.On("FindAppointmentByAppCode", "SLOT123").Return(appSlot, nil).Once()
		mockAppointmentRepo.On("FindAndLock", "SLOT123", mock.AnythingOfType("*gorm.DB")).Return(appSlot, nil).Once()
		// Transaction expectations
		sqlMock.ExpectBegin()
		mockBookingRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookingRepo).Once()
		mockBookingRepo.On("FindAndLockAvailableSlot", "SLOT123", slot.Date, slot.StartTime).Return(nil, fmt.Errorf("not found")).Once()
		sqlMock.ExpectRollback()

		// Act
		_, err := bookingService.BookAppointment(validReq, "")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "BOOKING_SLOT_UNAVAILABLE: no available slot found", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
		mockBookingRepo.AssertExpectations(t)
	})
}

func TestRefreshBookingStatuses(t *testing.T) {
	mockAppointmentRepo := new(repomocks.AppointmentRepository)
	mockBookingRepo := new(repomocks.BookingRepository)
	mockUserRepo := new(repomocks.UserRepository)
	mockBanListRepo := new(repomocks.BanListRepository)
	mockEventBus := new(MockEventBus)

	bookingService := services.NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, nil)

	now := time.Now()
	mockBookingRepo.On("MarkBookingsOngoing", mock.Anything, now).Return(int64(2), nil).Once()
	mockBookingRepo.On("MarkBookingsExpired", mock.Anything, now).Return(int64(1), nil).Once()

	summary, err := bookingService.RefreshBookingStatuses(context.Background(), now)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), summary.Ongoing)
	assert.Equal(t, int64(1), summary.Expired)
	mockBookingRepo.AssertExpectations(t)
}

func stubAppointmentWithTx(repo *repomocks.AppointmentRepository) {
	repo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(repo)
}

func stubBanListNotFound(repo *repomocks.BanListRepository) {
	repo.On("FindByUserAndEmail", mock.Anything, mock.Anything).
		Return((*entities.BanListEntry)(nil), repoerrors.NotFoundError("not found"))
}

func TestBookAppointmentAntiScalping(t *testing.T) {
	userID := uuid.New()
	user := &entities.User{ID: userID, Name: "Test User", Email: "test@example.com"}
	now := time.Now()
	slotStart := now.Add(1 * time.Hour)
	slotEnd := now.Add(2 * time.Hour)
	appStrict := &entities.Appointment{
		ID:                uuid.New(),
		AppCode:           "STRICT123",
		Type:              entities.Party,
		MaxAttendees:      10,
		AntiScalpingLevel: entities.ScalpingStrict,
		OwnerID:           uuid.New(),
		StartDate:         now,
		EndDate:           now,
		StartTime:         slotStart,
		EndTime:           slotEnd,
	}
	appStandard := &entities.Appointment{
		ID:                uuid.New(),
		AppCode:           "STD123",
		Type:              entities.Party,
		MaxAttendees:      10,
		AntiScalpingLevel: entities.ScalpingStandard,
		OwnerID:           uuid.New(),
		StartDate:         now,
		EndDate:           now,
		StartTime:         slotStart,
		EndTime:           slotEnd,
	}

	// Generate a valid token for success cases
	validDeviceID := "unique-device-123"
	validToken, err := middleware.GenerateDeviceToken(validDeviceID)
	assert.NoError(t, err)

	// Create a valid booking request with all required fields
	req := requests.BookingRequest{
		AppCode:       "STRICT123",
		AttendeeCount: 1,
		DeviceToken:   validToken,
		StartTime:     slotStart,
		EndTime:       slotEnd,
		Date:          now,
	}

	t.Run("Failure - Strict - Email already exists", func(t *testing.T) {
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockEventBus := new(MockEventBus)
		stubBanListNotFound(mockBanListRepo)

		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := services.NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, gormDB)

		mockAppointmentRepo.On("FindAppointmentByAppCode", "STRICT123").Return(appStrict, nil).Once()
		mockUserRepo.On("FindByID", userID.String()).Return(user, nil).Once()

		// Mock device check to PASS, so the code proceeds to the email check
		mockBookingRepo.On("FindActiveBookingByDevice", appStrict.ID, validDeviceID).Return(nil, repoerrors.NotFoundError("active booking not found for device")).Once()

		// Mock email check to FAIL
		mockBookingRepo.On("FindActiveBookingByEmail", appStrict.ID, user.Email).Return(&entities.Booking{}, nil).Once()

		_, err := bookingService.BookAppointment(req, userID.String())

		assert.Error(t, err)
		assert.Equal(t, "CONFLICT: this email has already been used to book for this appointment", err.Error())
		mockBookingRepo.AssertExpectations(t)
	})

	t.Run("Failure - Strict - Device already exists", func(t *testing.T) {
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		stubAppointmentWithTx(mockAppointmentRepo)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockEventBus := new(MockEventBus)
		stubBanListNotFound(mockBanListRepo)

		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := services.NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, gormDB)

		mockAppointmentRepo.On("FindAppointmentByAppCode", "STRICT123").Return(appStrict, nil).Once()
		mockUserRepo.On("FindByID", userID.String()).Return(user, nil).Once()

		// Mock device check to FAIL
		mockBookingRepo.On("FindActiveBookingByDevice", appStrict.ID, validDeviceID).Return(&entities.Booking{}, nil).Once()

		_, err := bookingService.BookAppointment(req, userID.String())

		assert.Error(t, err)
		assert.Equal(t, "CONFLICT: a booking has already been made from this device", err.Error())
		mockBookingRepo.AssertExpectations(t)
	})

	t.Run("Failure - Standard - Email already exists", func(t *testing.T) {
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		stubAppointmentWithTx(mockAppointmentRepo)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockEventBus := new(MockEventBus)
		stubBanListNotFound(mockBanListRepo)

		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := services.NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, gormDB)

		standardReq := requests.BookingRequest{
			AppCode:       "STD123",
			AttendeeCount: 1,
			StartTime:     now.Add(1 * time.Hour),
			EndTime:       now.Add(2 * time.Hour),
			Date:          now,
		}

		mockAppointmentRepo.On("FindAppointmentByAppCode", "STD123").Return(appStandard, nil).Once()
		mockUserRepo.On("FindByID", userID.String()).Return(user, nil).Once()
		// Mock email check to fail
		mockBookingRepo.On("FindActiveBookingByEmail", appStandard.ID, user.Email).Return(&entities.Booking{}, nil).Once()

		_, err := bookingService.BookAppointment(standardReq, userID.String())

		assert.Error(t, err)
		assert.Equal(t, "CONFLICT: this email has already been used to book for this appointment", err.Error())
		mockBookingRepo.AssertExpectations(t)
	})

	t.Run("Failure - Strict - Missing Device Token", func(t *testing.T) {
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		stubAppointmentWithTx(mockAppointmentRepo)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockEventBus := new(MockEventBus)
		stubBanListNotFound(mockBanListRepo)

		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := services.NewBookingService(nil, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, gormDB)

		missingTokenReq := requests.BookingRequest{
			AppCode:       "STRICT123",
			AttendeeCount: 1,
			StartTime:     now.Add(1 * time.Hour),
			EndTime:       now.Add(2 * time.Hour),
			Date:          now,
			DeviceToken:   "",
		}

		mockAppointmentRepo.On("FindAppointmentByAppCode", "STRICT123").Return(appStrict, nil).Once()
		mockUserRepo.On("FindByID", userID.String()).Return(user, nil).Once()

		_, err := bookingService.BookAppointment(missingTokenReq, userID.String())

		assert.Error(t, err)
		assert.Equal(t, "PRECONDITION_FAILED: device token is required for this appointment", err.Error())
	})
}

func TestUpdateBookingByCodeOngoingBlocked(t *testing.T) {
	mockAppointmentRepo := new(repomocks.AppointmentRepository)
	mockBookingRepo := new(repomocks.BookingRepository)
	mockUserRepo := new(repomocks.UserRepository)
	mockBanListRepo := new(repomocks.BanListRepository)
	mockEventBus := new(MockEventBus)

	bookingService := services.NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, nil)

	booking := &entities.Booking{Status: entities.BookingStatusOngoing, AppCode: "APP123"}
	mockBookingRepo.On("GetBookingByCode", "BK-ONGOING").Return(booking, nil).Once()

	_, err := bookingService.UpdateBookingByCode("BK-ONGOING", requests.BookingRequest{AppCode: "APP123"})

	assert.Error(t, err)
	assert.Equal(t, "CONFLICT: ongoing bookings cannot be rescheduled", err.Error())
	mockBookingRepo.AssertExpectations(t)
}

func TestUpdateBookingByCodeConfirmedResetsToPending(t *testing.T) {
	mockAppointmentRepo := new(repomocks.AppointmentRepository)
	mockBookingRepo := new(repomocks.BookingRepository)
	mockUserRepo := new(repomocks.UserRepository)
	mockBanListRepo := new(repomocks.BanListRepository)
	mockEventBus := new(MockEventBus)

	db, sqlMock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	bookingService := services.NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockEventBus, gormDB)

	oldDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	oldStart := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	newDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	newStart := time.Date(2026, 1, 1, 15, 45, 0, 0, time.UTC)
	newEnd := time.Date(2026, 1, 1, 16, 30, 0, 0, time.UTC)

	userID := uuid.New()
	booking := &entities.Booking{
		AppCode:       "APP123",
		Date:          oldDate,
		StartTime:     oldStart,
		EndTime:       oldStart.Add(45 * time.Minute),
		BookingCode:   "BK-OLD",
		Status:        entities.BookingStatusConfirmed,
		IsSlot:        true,
		AttendeeCount: 1,
		UserID:        &userID,
		Name:          "Test User",
		Email:         "test@example.com",
	}
	oldSlot := &entities.Booking{
		AppCode:     "APP123",
		Date:        oldDate,
		StartTime:   oldStart,
		EndTime:     oldStart.Add(45 * time.Minute),
		BookingCode: "BK-OLD",
		IsSlot:      true,
		Capacity:    1,
		SeatsBooked: 1,
		Status:      entities.BookingStatusConfirmed,
		Available:   false,
	}
	newSlot := &entities.Booking{
		AppCode:     "APP123",
		Date:        newDate,
		StartTime:   newStart,
		EndTime:     newEnd,
		BookingCode: "BK-NEW",
		IsSlot:      true,
		Capacity:    1,
		SeatsBooked: 0,
		Status:      entities.BookingStatusActive,
		Available:   true,
	}

	mockBookingRepo.On("GetBookingByCode", "BK-OLD").Return(booking, nil).Once()
	mockAppointmentRepo.On("FindAppointmentByAppCode", "APP123").Return(&entities.Appointment{AppCode: "APP123", Type: entities.Single}, nil).Once()
	sqlMock.ExpectBegin()
	mockBookingRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookingRepo).Once()
	mockBookingRepo.On("FindAndLockSlot", "APP123", oldDate, oldStart).Return(oldSlot, nil).Once()
	mockBookingRepo.On("FindAndLockSlot", "APP123", newDate, newStart).Return(newSlot, nil).Once()
	mockBookingRepo.On("Update", mock.AnythingOfType("*entities.Booking")).Return(nil).Times(3)
	sqlMock.ExpectCommit()

	mockEventBus.On("Publish", mock.Anything, mock.MatchedBy(func(event events.Event) bool {
		return event.Name == events.EventBookingUpdated
	})).Return(nil).Once()

	updated, err := bookingService.UpdateBookingByCode("BK-OLD", requests.BookingRequest{
		AppCode:       "APP123",
		Date:          newDate,
		StartTime:     newStart,
		EndTime:       newEnd,
		AttendeeCount: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, entities.BookingStatusPending, newSlot.Status)
	assert.Equal(t, "BK-OLD", newSlot.BookingCode)
	assert.Equal(t, updated, newSlot)
	mockBookingRepo.AssertExpectations(t)
	mockAppointmentRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestBanListService tests the ban list functionality
func TestBanListService(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"

	t.Run("Success - Add to ban list", func(t *testing.T) {
		mockBanListRepo := new(repomocks.BanListRepository)
		banService := services.NewBanListService(mockBanListRepo)

		mockBanListRepo.On("FindByUserAndEmail", userID, email).Return(nil, fmt.Errorf("not found")).Once()
		mockBanListRepo.On("Create", mock.AnythingOfType("*entities.BanListEntry")).Return(nil).Once()

		entry, err := banService.AddToBanList(userID, email)

		assert.NoError(t, err)
		assert.NotNil(t, entry)
		assert.Equal(t, userID, entry.UserID)
		assert.Equal(t, email, entry.BannedEmail)
		mockBanListRepo.AssertExpectations(t)
	})

	t.Run("Failure - Email already on ban list", func(t *testing.T) {
		mockBanListRepo := new(repomocks.BanListRepository)
		banService := services.NewBanListService(mockBanListRepo)

		existingEntry := &entities.BanListEntry{UserID: userID, BannedEmail: email}
		mockBanListRepo.On("FindByUserAndEmail", userID, email).Return(existingEntry, nil).Once()

		_, err := banService.AddToBanList(userID, email)

		assert.Error(t, err)
		assert.Equal(t, "CONFLICT: email already on ban list", err.Error())
		mockBanListRepo.AssertExpectations(t)
	})

	t.Run("Success - Remove from ban list", func(t *testing.T) {
		mockBanListRepo := new(repomocks.BanListRepository)
		banService := services.NewBanListService(mockBanListRepo)

		mockBanListRepo.On("Delete", userID, email).Return(nil).Once()

		err := banService.RemoveFromBanList(userID, email)

		assert.NoError(t, err)
		mockBanListRepo.AssertExpectations(t)
	})
}
