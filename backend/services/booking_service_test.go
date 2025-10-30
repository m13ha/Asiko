package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/middleware"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	notificationmocks "github.com/m13ha/appointment_master/notifications/mocks"
	repomocks "github.com/m13ha/appointment_master/repository/mocks"
	servicemocks "github.com/m13ha/appointment_master/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBookAppointment(t *testing.T) {
	userID := uuid.New()
	user := &entities.User{ID: userID, Name: "Test User", Email: "test@example.com"}
	appSlot := &entities.Appointment{ID: uuid.New(), AppCode: "SLOT123", Type: entities.Group, MaxAttendees: 5, AntiScalpingLevel: entities.ScalpingNone}
	slot := &entities.Booking{ID: uuid.New(), AppCode: "SLOT123", Available: true, Date: time.Now(), StartTime: time.Now().Add(time.Hour), EndTime: time.Now().Add(90 * time.Minute)}

    t.Run("Success - Book Slot", func(t *testing.T) {
        // Arrange
        mockAppointmentRepo := new(repomocks.AppointmentRepository)
        mockBookingRepo := new(repomocks.BookingRepository)
        mockUserRepo := new(repomocks.UserRepository)
        mockBanListRepo := new(repomocks.BanListRepository)
        mockNotificationService := new(notificationmocks.NotificationService)
        mockEventNotificationService := new(servicemocks.EventNotificationService)
        db, sqlMock, _ := sqlmock.New()
        gormDB, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

        bookingService := NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockNotificationService, mockEventNotificationService, gormDB)

		req := requests.BookingRequest{AppCode: "SLOT123", Name: "Test", Email: "test@test.com", Date: slot.Date, StartTime: slot.StartTime, EndTime: slot.EndTime, AttendeeCount: 2}

        mockAppointmentRepo.On("FindAppointmentByAppCode", "SLOT123").Return(appSlot, nil).Once()
        mockUserRepo.On("FindByID", userID.String()).Return(user, nil).Once()
        // Transaction expectations
        sqlMock.ExpectBegin()
        // Repository uses WithTx + FindAndLockAvailableSlot in transaction
        mockBookingRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookingRepo).Once()
        mockBookingRepo.On("FindAndLockAvailableSlot", "SLOT123", slot.Date, slot.StartTime).Return(slot, nil).Once()
        mockBookingRepo.On("Update", mock.AnythingOfType("*entities.Booking")).Return(nil).Once()
        sqlMock.ExpectCommit()
        mockNotificationService.On("SendBookingConfirmation", mock.AnythingOfType("*entities.Booking")).Return(nil).Once()
        mockEventNotificationService.On("CreateEventNotification", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
        mockBookingRepo.On("UpdateNotificationStatus", mock.AnythingOfType("uuid.UUID"), "sent", "email").Return(nil).Once()

		// Act
		booking, err := bookingService.BookAppointment(req, userID.String())

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, booking)
		mockAppointmentRepo.AssertExpectations(t)
		mockBookingRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
		mockEventNotificationService.AssertExpectations(t)
	})

	t.Run("Failure - Appointment not found", func(t *testing.T) {
		// Arrange
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockNotificationService := new(notificationmocks.NotificationService)
		mockEventNotificationService := new(servicemocks.EventNotificationService)
		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockNotificationService, mockEventNotificationService, gormDB)
		validReq := requests.BookingRequest{AppCode: "NOTFOUND", Name: "Guest User", Email: "guest@example.com", Date: time.Now(), StartTime: time.Now(), EndTime: time.Now(), AttendeeCount: 1}
		mockAppointmentRepo.On("FindAppointmentByAppCode", "NOTFOUND").Return(nil, fmt.Errorf("not found")).Once()

		// Act
		_, err := bookingService.BookAppointment(validReq, "")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "appointment not found", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
	})

    t.Run("Failure - Slot not available", func(t *testing.T) {
        // Arrange
        mockAppointmentRepo := new(repomocks.AppointmentRepository)
        mockBookingRepo := new(repomocks.BookingRepository)
        mockUserRepo := new(repomocks.UserRepository)
        mockBanListRepo := new(repomocks.BanListRepository)
        mockNotificationService := new(notificationmocks.NotificationService)
        mockEventNotificationService := new(servicemocks.EventNotificationService)
        db, sqlMock, _ := sqlmock.New()
        gormDB, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

		bookingService := NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockNotificationService, mockEventNotificationService, gormDB)
		validReq := requests.BookingRequest{AppCode: "SLOT123", Name: "Guest User", Email: "guest@example.com", Date: slot.Date, StartTime: slot.StartTime, EndTime: slot.EndTime, AttendeeCount: 1}
        mockAppointmentRepo.On("FindAppointmentByAppCode", "SLOT123").Return(appSlot, nil).Once()
        // Transaction expectations
        sqlMock.ExpectBegin()
        mockBookingRepo.On("WithTx", mock.AnythingOfType("*gorm.DB")).Return(mockBookingRepo).Once()
        mockBookingRepo.On("FindAndLockAvailableSlot", "SLOT123", slot.Date, slot.StartTime).Return(nil, fmt.Errorf("not found")).Once()
        sqlMock.ExpectRollback()

		// Act
		_, err := bookingService.BookAppointment(validReq, "")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "no available slot found", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
		mockBookingRepo.AssertExpectations(t)
	})
}

func TestBookAppointmentAntiScalping(t *testing.T) {
	userID := uuid.New()
	user := &entities.User{ID: userID, Name: "Test User", Email: "test@example.com"}
	appStrict := &entities.Appointment{ID: uuid.New(), AppCode: "STRICT123", Type: entities.Party, MaxAttendees: 10, AntiScalpingLevel: entities.ScalpingStrict, OwnerID: uuid.New()}
	appStandard := &entities.Appointment{ID: uuid.New(), AppCode: "STD123", Type: entities.Party, MaxAttendees: 10, AntiScalpingLevel: entities.ScalpingStandard, OwnerID: uuid.New()}

	// Generate a valid token for success cases
	validDeviceID := "unique-device-123"
	validToken, err := middleware.GenerateDeviceToken(validDeviceID)
	assert.NoError(t, err)

	// Create a valid booking request with all required fields
	now := time.Now()
	req := requests.BookingRequest{
		AppCode:       "STRICT123",
		AttendeeCount: 1,
		DeviceToken:   validToken,
		StartTime:     now.Add(1 * time.Hour),
		EndTime:       now.Add(2 * time.Hour),
		Date:          now,
	}

	t.Run("Failure - Strict - Email already exists", func(t *testing.T) {
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockNotificationService := new(notificationmocks.NotificationService)
		mockEventNotificationService := new(servicemocks.EventNotificationService)

		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockNotificationService, mockEventNotificationService, gormDB)

		mockAppointmentRepo.On("FindAppointmentByAppCode", "STRICT123").Return(appStrict, nil).Once()
		mockUserRepo.On("FindByID", userID.String()).Return(user, nil).Once()

		// Mock device check to PASS, so the code proceeds to the email check
		mockBookingRepo.On("FindActiveBookingByDevice", appStrict.ID, validDeviceID).Return(nil, gorm.ErrRecordNotFound).Once()

		// Mock email check to FAIL
		mockBookingRepo.On("FindActiveBookingByEmail", appStrict.ID, user.Email).Return(&entities.Booking{}, nil).Once()

		_, err := bookingService.BookAppointment(req, userID.String())

		assert.Error(t, err)
		assert.Equal(t, "this email has already been used to book for this appointment", err.Error())
		mockBookingRepo.AssertExpectations(t)
	})

	t.Run("Failure - Strict - Device already exists", func(t *testing.T) {
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockNotificationService := new(notificationmocks.NotificationService)
		mockEventNotificationService := new(servicemocks.EventNotificationService)

		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockNotificationService, mockEventNotificationService, gormDB)

		mockAppointmentRepo.On("FindAppointmentByAppCode", "STRICT123").Return(appStrict, nil).Once()
		mockUserRepo.On("FindByID", userID.String()).Return(user, nil).Once()

		// Mock device check to FAIL
		mockBookingRepo.On("FindActiveBookingByDevice", appStrict.ID, validDeviceID).Return(&entities.Booking{}, nil).Once()

		_, err := bookingService.BookAppointment(req, userID.String())

		assert.Error(t, err)
		assert.Equal(t, "a booking has already been made from this device", err.Error())
		mockBookingRepo.AssertExpectations(t)
	})

	t.Run("Failure - Standard - Email already exists", func(t *testing.T) {
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		mockBookingRepo := new(repomocks.BookingRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockNotificationService := new(notificationmocks.NotificationService)
		mockEventNotificationService := new(servicemocks.EventNotificationService)

		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := NewBookingService(mockBookingRepo, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockNotificationService, mockEventNotificationService, gormDB)

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
		assert.Equal(t, "this email has already been used to book for this appointment", err.Error())
		mockBookingRepo.AssertExpectations(t)
	})

	t.Run("Failure - Strict - Missing Device Token", func(t *testing.T) {
		mockAppointmentRepo := new(repomocks.AppointmentRepository)
		mockUserRepo := new(repomocks.UserRepository)
		mockBanListRepo := new(repomocks.BanListRepository)
		mockNotificationService := new(notificationmocks.NotificationService)
		mockEventNotificationService := new(servicemocks.EventNotificationService)

		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})

		bookingService := NewBookingService(nil, mockAppointmentRepo, mockUserRepo, mockBanListRepo, mockNotificationService, mockEventNotificationService, gormDB)

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
		assert.Equal(t, "device token is required for this appointment", err.Error())
	})
}

// TestBanListService tests the ban list functionality
func TestBanListService(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"

	t.Run("Success - Add to ban list", func(t *testing.T) {
		mockBanListRepo := new(repomocks.BanListRepository)
		banService := NewBanListService(mockBanListRepo)

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
		banService := NewBanListService(mockBanListRepo)

		existingEntry := &entities.BanListEntry{UserID: userID, BannedEmail: email}
		mockBanListRepo.On("FindByUserAndEmail", userID, email).Return(existingEntry, nil).Once()

		_, err := banService.AddToBanList(userID, email)

		assert.Error(t, err)
		assert.Equal(t, "email already on ban list", err.Error())
		mockBanListRepo.AssertExpectations(t)
	})

	t.Run("Success - Remove from ban list", func(t *testing.T) {
		mockBanListRepo := new(repomocks.BanListRepository)
		banService := NewBanListService(mockBanListRepo)

		mockBanListRepo.On("Delete", userID, email).Return(nil).Once()

		err := banService.RemoveFromBanList(userID, email)

		assert.NoError(t, err)
		mockBanListRepo.AssertExpectations(t)
	})
}
