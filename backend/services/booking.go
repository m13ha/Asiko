package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	myerrors "github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/middleware"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/repository"
	"github.com/m13ha/appointment_master/utils"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type bookingServiceImpl struct {
	bookingRepo     repository.BookingRepository
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository
	banListRepo     repository.BanListRepository
	db              *gorm.DB
}

func NewBookingService(bookingRepo repository.BookingRepository, appointmentRepo repository.AppointmentRepository, userRepo repository.UserRepository, banListRepo repository.BanListRepository, db *gorm.DB) BookingService {
	return &bookingServiceImpl{bookingRepo: bookingRepo, appointmentRepo: appointmentRepo, userRepo: userRepo, banListRepo: banListRepo, db: db}
}

// performAntiScalpingChecks runs validation based on the appointment's settings.
// It returns the trusted device ID (if applicable) or an error if a check fails.
func (s *bookingServiceImpl) performAntiScalpingChecks(appointment *entities.Appointment, req requests.BookingRequest, bookingEmail string) (string, error) {
	level := appointment.AntiScalpingLevel
	if level == entities.ScalpingNone {
		return "", nil // No checks needed
	}

	var trustedDeviceID string

	// Strict check (Device ID)
	if level == entities.ScalpingStrict {
		if req.DeviceToken == "" {
			return "", myerrors.NewUserError("device token is required for this appointment")
		}
		validatedDeviceID, err := middleware.ValidateDeviceToken(req.DeviceToken)
		if err != nil {
			return "", myerrors.NewUserError(fmt.Sprintf("invalid device token: %v", err))
		}
		trustedDeviceID = validatedDeviceID

		// Check if device has already booked
		if _, err := s.bookingRepo.FindActiveBookingByDevice(appointment.ID, trustedDeviceID); err == nil {
			return "", myerrors.NewUserError("a booking has already been made from this device")
		}
	}

	// Standard check (Email) - runs for both 'standard' and 'strict'
	if level == entities.ScalpingStandard || level == entities.ScalpingStrict {
		if _, err := s.bookingRepo.FindActiveBookingByEmail(appointment.ID, bookingEmail); err == nil {
			return "", myerrors.NewUserError("this email has already been used to book for this appointment")
		}
	}

	return trustedDeviceID, nil
}

// BookAppointment handles booking for both registered users and guests
func (s *bookingServiceImpl) BookAppointment(req requests.BookingRequest, userIDStr string) (*entities.Booking, error) {
	// --- 1. Basic Validation ---
	if userIDStr == "" {
		if err := req.Validate(); err != nil {
			return nil, err
		}
	} else {
		if err := utils.Validate(req); err != nil {
			return nil, myerrors.NewUserError("invalid booking data: " + err.Error())
		}
	}

	// --- 2. Fetch Appointment ---
	appointment, err := s.appointmentRepo.FindAppointmentByAppCode(req.AppCode)
	if err != nil {
		return nil, myerrors.NewUserError("appointment not found")
	}

	// --- 3. Get Booker's Info ---
	var user *entities.User
	var bookingEmail string
	if userIDStr != "" {
		user, err = s.userRepo.FindByID(userIDStr)
		if err != nil {
			return nil, myerrors.NewUserError("user not found")
		}
		bookingEmail = user.Email
	} else {
		bookingEmail = utils.NormalizeEmail(req.Email)
	}

	// --- 4. Anti-Scalping Checks ---
	trustedDeviceID, err := s.performAntiScalpingChecks(appointment, req, bookingEmail)
	if err != nil {
		return nil, err
	}

	// --- 5. Proceed with Booking ---
	if appointment.Type == entities.Party {
		return s.bookPartyAppointment(req, user, appointment, trustedDeviceID)
	}
	return s.bookSlotAppointment(req, user, appointment, trustedDeviceID)
}

func (s *bookingServiceImpl) bookPartyAppointment(req requests.BookingRequest, user *entities.User, appointment *entities.Appointment, deviceID string) (*entities.Booking, error) {
	var booking *entities.Booking
	err := s.db.Transaction(func(tx *gorm.DB) error {
		appRepo := s.appointmentRepo.WithTx(tx)
		bookRepo := s.bookingRepo.WithTx(tx)

		lockedAppointment, err := appRepo.FindAndLock(req.AppCode, tx)
		if err != nil {
			return myerrors.NewUserError("appointment not found")
		}

		if lockedAppointment.AttendeesBooked+req.AttendeeCount > lockedAppointment.MaxAttendees {
			return myerrors.NewUserError("not enough capacity for this party")
		}

		booking = &entities.Booking{
			AppointmentID: lockedAppointment.ID,
			AppCode:       lockedAppointment.AppCode,
			Date:          req.Date,
			StartTime:     req.StartTime,
			EndTime:       req.EndTime,
			Available:     false,
			AttendeeCount: req.AttendeeCount,
			Description:   req.Description,
			Status:        "active",
			DeviceID:      deviceID,
		}

		if user != nil {
			booking.UserID = &user.ID
			booking.Name = user.Name
			booking.Email = user.Email
			booking.Phone = user.PhoneNumber
		} else {
			booking.Name = req.Name
			booking.Email = utils.NormalizeEmail(req.Email)
			booking.Phone = req.Phone
		}

		if err := bookRepo.Create(booking); err != nil {
			return err
		}

		lockedAppointment.AttendeesBooked += req.AttendeeCount
		return appRepo.Update(lockedAppointment)
	})

	return booking, err
}

func (s *bookingServiceImpl) bookSlotAppointment(req requests.BookingRequest, user *entities.User, appointment *entities.Appointment, deviceID string) (*entities.Booking, error) {
	slot, err := s.bookingRepo.FindAvailableSlot(req.AppCode, req.Date, req.StartTime)
	if err != nil {
		return nil, myerrors.NewUserError("no available slot found")
	}

	// Populate user info
	if user != nil {
		slot.UserID = &user.ID
		slot.Name = user.Name
		slot.Email = user.Email
		slot.Phone = user.PhoneNumber
	} else {
		slot.Name = req.Name
		slot.Email = utils.NormalizeEmail(req.Email)
		slot.Phone = req.Phone
	}

	// Populate booking details
	if appointment.Type == entities.Group {
		if req.AttendeeCount > appointment.MaxAttendees {
			return nil, myerrors.NewUserError("attendee count exceeds maximum allowed")
		}
		slot.AttendeeCount = req.AttendeeCount
	} else {
		slot.AttendeeCount = 1
	}
	slot.Description = req.Description
	slot.DeviceID = deviceID
	slot.Available = false

	if err := s.bookingRepo.Update(slot); err != nil {
		log.Printf("[bookSlot] DB error: %v", err)
		return nil, fmt.Errorf("internal error")
	}

	return slot, nil
}

// BookRegisteredUserAppointment is a wrapper for backward compatibility
func (s *bookingServiceImpl) BookRegisteredUserAppointment(req requests.BookingRequest, userIDStr string) (*entities.Booking, error) {
	return s.BookAppointment(req, userIDStr)
}

// BookGuestAppointment is a wrapper for backward compatibility
func (s *bookingServiceImpl) BookGuestAppointment(req requests.BookingRequest) (*entities.Booking, error) {
	return s.BookAppointment(req, "")
}

// GetAllBookingsForAppointment returns all bookings for a specific appointment with pagination
func (s *bookingServiceImpl) GetAllBookingsForAppointment(ctx context.Context, appcode string) (paginate.Page, error) {
	return s.bookingRepo.GetBookingsByAppCode(ctx, appcode, false), nil
}

// GetUserBookings returns all bookings for a specific user with pagination
func (s *bookingServiceImpl) GetUserBookings(ctx context.Context, userID string) (paginate.Page, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return paginate.Page{}, myerrors.NewUserError("Invalid user ID.")
	}
	return s.bookingRepo.GetBookingsByUserID(ctx, uid), nil
}

// GetAvailableSlots returns all available slots for an appointment with pagination
func (s *bookingServiceImpl) GetAvailableSlots(ctx context.Context, appcode string) (paginate.Page, error) {
	return s.bookingRepo.GetAvailableSlots(ctx, appcode), nil
}

// GetAvailableSlotsByDay returns available slots for an appointment on a specific day with pagination
func (s *bookingServiceImpl) GetAvailableSlotsByDay(ctx context.Context, appcode string, dateStr string) (paginate.Page, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return paginate.Page{}, myerrors.NewUserError("Invalid date format. Use YYYY-MM-DD.")
	}
	return s.bookingRepo.GetAvailableSlotsByDay(ctx, appcode, parsedDate), nil
}

// GetBookingByCode retrieves a booking by its permanent booking_code
func (s *bookingServiceImpl) GetBookingByCode(bookingCode string) (*entities.Booking, error) {
	booking, err := s.bookingRepo.GetBookingByCode(bookingCode)
	if err != nil {
		return nil, myerrors.NewUserError("Booking not found.")
	}
	return booking, nil
}

// UpdateBookingByCode allows rescheduling a booking if the new slot is available
func (s *bookingServiceImpl) UpdateBookingByCode(bookingCode string, req requests.BookingRequest) (*entities.Booking, error) {
	booking, err := s.GetBookingByCode(bookingCode)
	if err != nil {
		return nil, err
	}
	_, err = s.bookingRepo.FindAvailableSlot(req.AppCode, req.Date, req.StartTime)
	if err != nil {
		return nil, myerrors.NewUserError("Requested slot is not available.")
	}
	booking.Available = true
	s.bookingRepo.Update(booking)
	booking.Date = req.Date
	booking.StartTime = req.StartTime
	booking.EndTime = req.EndTime
	booking.AttendeeCount = req.AttendeeCount
	booking.Description = req.Description
	booking.Available = false
	booking.Status = "active"
	if err := s.bookingRepo.Update(booking); err != nil {
		log.Printf("[UpdateBookingByCode] DB error: %v", err)
		return nil, fmt.Errorf("internal error")
	}
	return booking, nil
}

// CancelBookingByCode cancels a booking by booking_code
func (s *bookingServiceImpl) CancelBookingByCode(bookingCode string) (*entities.Booking, error) {
	booking, err := s.GetBookingByCode(bookingCode)
	if err != nil {
		return nil, err
	}

	appointment, err := s.appointmentRepo.FindAppointmentByAppCode(booking.AppCode)
	if err != nil {
		return nil, myerrors.NewUserError("Appointment not found.")
	}

	if appointment.Type == entities.Party {
		err := s.db.Transaction(func(tx *gorm.DB) error {
			appRepo := s.appointmentRepo.WithTx(tx)
			bookRepo := s.bookingRepo.WithTx(tx)

			lockedAppointment, err := appRepo.FindAndLock(appointment.AppCode, tx)
			if err != nil {
				return err
			}

			lockedAppointment.AttendeesBooked -= booking.AttendeeCount
			if err := appRepo.Update(lockedAppointment); err != nil {
				return err
			}

			booking.Status = "cancelled"
			if err := bookRepo.Update(booking); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Printf("[CancelBookingByCode] DB error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
	} else {
		booking.Available = true
		booking.Status = "cancelled"
		if err := s.bookingRepo.Update(booking); err != nil {
			log.Printf("[CancelBookingByCode] DB error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
	}

	return booking, nil
}

func (s *bookingServiceImpl) RejectBooking(bookingCode string, ownerID uuid.UUID) (*entities.Booking, error) {
	booking, err := s.GetBookingByCode(bookingCode)
	if err != nil {
		return nil, err
	}

	appointment, err := s.appointmentRepo.FindAppointmentByAppCode(booking.AppCode)
	if err != nil {
		return nil, myerrors.NewUserError("Appointment not found.")
	}

	if appointment.OwnerID != ownerID {
		return nil, myerrors.NewUserError("you are not the owner of this appointment")
	}

	if appointment.Type == entities.Party {
		err := s.db.Transaction(func(tx *gorm.DB) error {
			appRepo := s.appointmentRepo.WithTx(tx)
			bookRepo := s.bookingRepo.WithTx(tx)

			lockedAppointment, err := appRepo.FindAndLock(appointment.AppCode, tx)
			if err != nil {
				return err
			}

			lockedAppointment.AttendeesBooked -= booking.AttendeeCount
			if err := appRepo.Update(lockedAppointment); err != nil {
				return err
			}

			booking.Status = "rejected"
			if err := bookRepo.Update(booking); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Printf("[RejectBooking] DB error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
	} else {
		booking.Available = true
		booking.Status = "rejected"
		if err := s.bookingRepo.Update(booking); err != nil {
			log.Printf("[RejectBooking] DB error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
	}

	return booking, nil
}
