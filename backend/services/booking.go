package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	serviceerrors "github.com/m13ha/asiko/errors/serviceerrors"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	"github.com/m13ha/asiko/notifications"
	"github.com/m13ha/asiko/repository"
	"github.com/m13ha/asiko/utils"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type bookingServiceImpl struct {
	bookingRepo              repository.BookingRepository
	appointmentRepo          repository.AppointmentRepository
	userRepo                 repository.UserRepository
	banListRepo              repository.BanListRepository
	notificationService      notifications.NotificationService
	eventNotificationService EventNotificationService
	db                       *gorm.DB
}

type BookingStatusRefreshSummary struct {
	Ongoing int64
	Expired int64
}

func NewBookingService(bookingRepo repository.BookingRepository, appointmentRepo repository.AppointmentRepository, userRepo repository.UserRepository, banListRepo repository.BanListRepository, notificationService notifications.NotificationService, eventNotificationService EventNotificationService, db *gorm.DB) BookingService {
	return &bookingServiceImpl{bookingRepo: bookingRepo, appointmentRepo: appointmentRepo, userRepo: userRepo, banListRepo: banListRepo, notificationService: notificationService, eventNotificationService: eventNotificationService, db: db}
}

func ensureAppointmentIsBookable(appointment *entities.Appointment) error {
	switch appointment.Status {
	case entities.AppointmentStatusCanceled:
		return serviceerrors.ConflictError("appointment has been canceled")
	case entities.AppointmentStatusCompleted, entities.AppointmentStatusExpired:
		return serviceerrors.ConflictError("appointment is no longer active")
	default:
		return nil
	}
}

func normalizeSlotState(slot *entities.Booking) {
	if slot.Capacity < 1 {
		slot.Capacity = 1
	}
	if slot.SeatsBooked < 0 {
		slot.SeatsBooked = 0
	}
	if slot.SeatsBooked >= slot.Capacity {
		slot.Available = false
	} else {
		slot.Available = true
	}
	remaining := slot.Capacity - slot.SeatsBooked
	if remaining < 0 {
		remaining = 0
	}
	slot.AttendeeCount = remaining
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
			return "", serviceerrors.PreconditionFailedError("device token is required for this appointment")
		}
		validatedDeviceID, err := middleware.ValidateDeviceToken(req.DeviceToken)
		if err != nil {
			return "", serviceerrors.ValidationError(fmt.Sprintf("invalid device token: %v", err))
		}
		trustedDeviceID = validatedDeviceID

		// Check if device has already booked
		if _, err := s.bookingRepo.FindActiveBookingByDevice(appointment.ID, trustedDeviceID); err == nil {
			return "", serviceerrors.ConflictError("a booking has already been made from this device")
		}
	}

	// Standard check (Email) - runs for both 'standard' and 'strict'
	if level == entities.ScalpingStandard || level == entities.ScalpingStrict {
		if _, err := s.bookingRepo.FindActiveBookingByEmail(appointment.ID, bookingEmail); err == nil {
			return "", serviceerrors.ConflictError("this email has already been used to book for this appointment")
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
			return nil, serviceerrors.UserError("invalid booking data: " + err.Error())
		}
	}

	// --- 2. Fetch Appointment ---
	appointment, err := s.appointmentRepo.FindAppointmentByAppCode(req.AppCode)
	if err != nil {
		return nil, err
	}
	if err := ensureAppointmentIsBookable(appointment); err != nil {
		return nil, err
	}

	// --- 3. Get Booker's Info ---
	var user *entities.User
	var bookingEmail string
	if userIDStr != "" {
		user, err = s.userRepo.FindByID(userIDStr)
		if err != nil {
			return nil, err
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
			return err
		}
		if err := ensureAppointmentIsBookable(lockedAppointment); err != nil {
			return err
		}

		if lockedAppointment.AttendeesBooked+req.AttendeeCount > lockedAppointment.MaxAttendees {
			return serviceerrors.BookingCapacityExceededError("not enough capacity for this party")
		}

		booking = &entities.Booking{
			AppointmentID: lockedAppointment.ID,
			AppCode:       lockedAppointment.AppCode,
			Date:          req.Date,
			StartTime:     req.StartTime,
			EndTime:       req.EndTime,
			Available:     false,
			IsSlot:        false,
			Capacity:      req.AttendeeCount,
			SeatsBooked:   req.AttendeeCount,
			AttendeeCount: req.AttendeeCount,
			Description:   req.Description,
			Status:        "active",
			DeviceID:      deviceID,
		}

		if user != nil {
			booking.UserID = &user.ID
			booking.Name = user.Name
			booking.Email = user.Email
			if user.PhoneNumber != nil {
				booking.Phone = *user.PhoneNumber
			} else {
				booking.Phone = ""
			}
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

	if err == nil {
		go func() {
			if err := s.notificationService.SendBookingConfirmation(booking); err != nil {
				s.bookingRepo.UpdateNotificationStatus(booking.ID, "failed", "email")
			} else {
				s.bookingRepo.UpdateNotificationStatus(booking.ID, "sent", "email")
			}
		}()

		message := fmt.Sprintf("New booking by %s for your appointment %s.", booking.Name, appointment.Title)
		s.eventNotificationService.CreateEventNotification(appointment.OwnerID, "BOOKING_CREATED", message, booking.ID)
	}

	return booking, err
}

func (s *bookingServiceImpl) bookSlotAppointment(req requests.BookingRequest, user *entities.User, appointment *entities.Appointment, deviceID string) (*entities.Booking, error) {
	if appointment.Type == entities.Group {
		var reservation *entities.Booking
		err := s.db.Transaction(func(tx *gorm.DB) error {
			appRepo := s.appointmentRepo.WithTx(tx)
			bookRepo := s.bookingRepo.WithTx(tx)

			lockedAppointment, err := appRepo.FindAndLock(req.AppCode, tx)
			if err != nil {
				return err
			}
			if err := ensureAppointmentIsBookable(lockedAppointment); err != nil {
				return err
			}

			lockedSlot, err := bookRepo.FindAndLockAvailableSlot(req.AppCode, req.Date, req.StartTime)
			if err != nil {
				return serviceerrors.BookingSlotUnavailableError("no available slot found")
			}

			remaining := lockedSlot.Capacity - lockedSlot.SeatsBooked
			if remaining <= 0 || req.AttendeeCount > remaining {
				return serviceerrors.BookingCapacityExceededError("not enough capacity for this slot")
			}

			if req.AttendeeCount > appointment.MaxAttendees {
				return serviceerrors.ValidationError("attendee count exceeds maximum allowed")
			}

			reservation = &entities.Booking{
				AppointmentID: lockedSlot.AppointmentID,
				AppCode:       lockedSlot.AppCode,
				Date:          lockedSlot.Date,
				StartTime:     lockedSlot.StartTime,
				EndTime:       lockedSlot.EndTime,
				Available:     false,
				IsSlot:        false,
				Capacity:      req.AttendeeCount,
				SeatsBooked:   req.AttendeeCount,
				AttendeeCount: req.AttendeeCount,
				Description:   req.Description,
				DeviceID:      deviceID,
			}

			if user != nil {
				reservation.UserID = &user.ID
				reservation.Name = user.Name
				reservation.Email = user.Email
				if user.PhoneNumber != nil {
					reservation.Phone = *user.PhoneNumber
				}
			} else {
				reservation.Name = req.Name
				reservation.Email = utils.NormalizeEmail(req.Email)
				reservation.Phone = req.Phone
			}

			if err := bookRepo.Create(reservation); err != nil {
				log.Printf("[bookSlot] failed to create reservation: %v", err)
				return serviceerrors.FromError(err)
			}

			lockedSlot.SeatsBooked += req.AttendeeCount
			normalizeSlotState(lockedSlot)

			if err := bookRepo.Update(lockedSlot); err != nil {
				log.Printf("[bookSlot] failed to update slot: %v", err)
				return serviceerrors.FromError(err)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		go func() {
			if err := s.notificationService.SendBookingConfirmation(reservation); err != nil {
				s.bookingRepo.UpdateNotificationStatus(reservation.ID, "failed", "email")
			} else {
				s.bookingRepo.UpdateNotificationStatus(reservation.ID, "sent", "email")
			}
		}()

		message := fmt.Sprintf("New booking by %s for your appointment %s.", reservation.Name, appointment.Title)
		s.eventNotificationService.CreateEventNotification(appointment.OwnerID, "BOOKING_CREATED", message, reservation.ID)

		return reservation, nil
	}

	// Fallback to single-slot behaviour for other appointment types
	var slot *entities.Booking
	err := s.db.Transaction(func(tx *gorm.DB) error {
		appRepo := s.appointmentRepo.WithTx(tx)
		bookRepo := s.bookingRepo.WithTx(tx)

		lockedAppointment, err := appRepo.FindAndLock(req.AppCode, tx)
		if err != nil {
			return err
		}
		if err := ensureAppointmentIsBookable(lockedAppointment); err != nil {
			return err
		}

		lockedSlot, err := bookRepo.FindAndLockAvailableSlot(req.AppCode, req.Date, req.StartTime)
		if err != nil {
			return serviceerrors.BookingSlotUnavailableError("no available slot found")
		}

		// Populate user info
		if user != nil {
			lockedSlot.UserID = &user.ID
			lockedSlot.Name = user.Name
			lockedSlot.Email = user.Email
			if user.PhoneNumber != nil {
				lockedSlot.Phone = *user.PhoneNumber
			} else {
				lockedSlot.Phone = ""
			}
		} else {
			lockedSlot.Name = req.Name
			lockedSlot.Email = utils.NormalizeEmail(req.Email)
			lockedSlot.Phone = req.Phone
		}

		lockedSlot.SeatsBooked = lockedSlot.Capacity
		lockedSlot.Description = req.Description
		lockedSlot.DeviceID = deviceID
		normalizeSlotState(lockedSlot)

		if err := bookRepo.Update(lockedSlot); err != nil {
			log.Printf("[bookSlot] DB error: %v", err)
			return serviceerrors.FromError(err)
		}

		slot = lockedSlot
		return nil
	})

	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.notificationService.SendBookingConfirmation(slot); err != nil {
			s.bookingRepo.UpdateNotificationStatus(slot.ID, "failed", "email")
		} else {
			s.bookingRepo.UpdateNotificationStatus(slot.ID, "sent", "email")
		}
	}()

	message := fmt.Sprintf("New booking by %s for your appointment %s.", slot.Name, appointment.Title)
	s.eventNotificationService.CreateEventNotification(appointment.OwnerID, "BOOKING_CREATED", message, slot.ID)

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
func (s *bookingServiceImpl) GetAllBookingsForAppointment(ctx context.Context, req *http.Request, appcode string) (paginate.Page, error) {
	return s.bookingRepo.GetBookingsByAppCode(ctx, req, appcode, false), nil
}

// GetUserBookings returns all bookings for a specific user with pagination
func (s *bookingServiceImpl) GetUserBookings(ctx context.Context, req *http.Request, userID string, statuses []string) (paginate.Page, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return paginate.Page{}, serviceerrors.ValidationError("Invalid user ID.")
	}
	return s.bookingRepo.GetBookingsByUserID(ctx, req, uid, statuses), nil
}

func (s *bookingServiceImpl) RefreshBookingStatuses(ctx context.Context, now time.Time) (BookingStatusRefreshSummary, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if now.IsZero() {
		now = time.Now()
	}

	var summary BookingStatusRefreshSummary
	updated, err := s.bookingRepo.MarkBookingsOngoing(ctx, now)
	if err != nil {
		return summary, err
	}
	summary.Ongoing = updated

	updated, err = s.bookingRepo.MarkBookingsExpired(ctx, now)
	if err != nil {
		return summary, err
	}
	summary.Expired = updated

	return summary, nil
}

// GetAvailableSlots returns all available slots for an appointment with pagination
func (s *bookingServiceImpl) GetAvailableSlots(req *http.Request, appcode string) (paginate.Page, error) {
	ctx := context.Background()
	if req != nil {
		ctx = req.Context()
	}
	return s.bookingRepo.GetAvailableSlots(ctx, req, appcode), nil
}

// GetAvailableSlotsByDay returns available slots for an appointment on a specific day with pagination
func (s *bookingServiceImpl) GetAvailableSlotsByDay(req *http.Request, appcode string, dateStr string) (paginate.Page, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return paginate.Page{}, serviceerrors.ValidationError("Invalid date format. Use YYYY-MM-DD.")
	}
	ctx := context.Background()
	if req != nil {
		ctx = req.Context()
	}
	return s.bookingRepo.GetAvailableSlotsByDay(ctx, req, appcode, parsedDate), nil
}

// GetAvailableDates returns a list of distinct dates with available slots
func (s *bookingServiceImpl) GetAvailableDates(ctx context.Context, appcode string) ([]string, error) {
	dates, err := s.bookingRepo.GetAvailableDates(ctx, appcode)
	if err != nil {
		return nil, serviceerrors.FromError(err)
	}

	dateStrings := make([]string, len(dates))
	for i, date := range dates {
		dateStrings[i] = date.Format("2006-01-02")
	}
	return dateStrings, nil
}

// GetBookingByCode retrieves a booking by its permanent booking_code
func (s *bookingServiceImpl) GetBookingByCode(bookingCode string) (*entities.Booking, error) {
	booking, err := s.bookingRepo.GetBookingByCode(bookingCode)
	if err != nil {
		return nil, serviceerrors.FromError(err)
	}
	return booking, nil
}

// UpdateBookingByCode allows rescheduling a booking if the new slot is available
func (s *bookingServiceImpl) UpdateBookingByCode(bookingCode string, req requests.BookingRequest) (*entities.Booking, error) {
	booking, err := s.GetBookingByCode(bookingCode)
	if err != nil {
		return nil, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		bookRepo := s.bookingRepo.WithTx(tx)
		capacityErr := serviceerrors.BookingCapacityExceededError("not enough capacity for this slot")
		sameSlot := booking.AppCode == req.AppCode && booking.Date.Equal(req.Date) && booking.StartTime.Equal(req.StartTime)

		if sameSlot {
			slot, slotErr := bookRepo.FindAndLockSlot(req.AppCode, req.Date, req.StartTime)
			if slotErr != nil {
				return slotErr
			}
			delta := req.AttendeeCount - booking.AttendeeCount
			if delta > 0 && slot.SeatsBooked+delta > slot.Capacity {
				return capacityErr
			}
			slot.SeatsBooked += delta
			normalizeSlotState(slot)
			if updateErr := bookRepo.Update(slot); updateErr != nil {
				return updateErr
			}
		} else {
			oldSlot, slotErr := bookRepo.FindAndLockSlot(booking.AppCode, booking.Date, booking.StartTime)
			if slotErr != nil {
				return slotErr
			}
			oldSlot.SeatsBooked -= booking.AttendeeCount
			normalizeSlotState(oldSlot)
			if updateErr := bookRepo.Update(oldSlot); updateErr != nil {
				return updateErr
			}

			newSlot, slotErr := bookRepo.FindAndLockSlot(req.AppCode, req.Date, req.StartTime)
			if slotErr != nil {
				return slotErr
			}
			if newSlot.SeatsBooked+req.AttendeeCount > newSlot.Capacity {
				return capacityErr
			}
			newSlot.SeatsBooked += req.AttendeeCount
			normalizeSlotState(newSlot)
			if updateErr := bookRepo.Update(newSlot); updateErr != nil {
				return updateErr
			}
		}

		booking.AppCode = req.AppCode
		booking.Date = req.Date
		booking.StartTime = req.StartTime
		booking.EndTime = req.EndTime
		booking.AttendeeCount = req.AttendeeCount
		booking.Description = req.Description
		booking.Available = false
		booking.Status = "active"

		if updateErr := bookRepo.Update(booking); updateErr != nil {
			return updateErr
		}

		return nil
	})

	if err != nil {
		return nil, serviceerrors.FromError(err)
	}

	appointment, _ := s.appointmentRepo.FindAppointmentByAppCode(booking.AppCode)
	if appointment != nil {
		message := fmt.Sprintf("Booking %s was updated.", booking.BookingCode)
		s.eventNotificationService.CreateEventNotification(appointment.OwnerID, "BOOKING_UPDATED", message, booking.ID)
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
		return nil, serviceerrors.FromError(err)
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
			return nil, serviceerrors.FromError(err)
		}
	} else {
		err = s.db.Transaction(func(tx *gorm.DB) error {
			bookRepo := s.bookingRepo.WithTx(tx)

			slot, slotErr := bookRepo.FindAndLockSlot(booking.AppCode, booking.Date, booking.StartTime)
			if slotErr != nil {
				return slotErr
			}

			slot.SeatsBooked -= booking.AttendeeCount
			normalizeSlotState(slot)
			if updateErr := bookRepo.Update(slot); updateErr != nil {
				return updateErr
			}

			booking.Available = true
			booking.Status = "cancelled"
			if updateErr := bookRepo.Update(booking); updateErr != nil {
				return updateErr
			}

			return nil
		})

		if err != nil {
			log.Printf("[CancelBookingByCode] DB error: %v", err)
			return nil, serviceerrors.FromError(err)
		}
	}

	if err == nil {
		go func() {
			if err := s.notificationService.SendBookingCancellation(booking); err != nil {
				s.bookingRepo.UpdateNotificationStatus(booking.ID, "failed", "email")
			} else {
				s.bookingRepo.UpdateNotificationStatus(booking.ID, "sent", "email")
			}
		}()

		message := fmt.Sprintf("Booking %s was cancelled.", booking.BookingCode)
		s.eventNotificationService.CreateEventNotification(appointment.OwnerID, "BOOKING_CANCELLED", message, booking.ID)
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
		return nil, serviceerrors.FromError(err)
	}

	if appointment.OwnerID != ownerID {
		return nil, serviceerrors.ForbiddenError("you are not the owner of this appointment")
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
			return nil, serviceerrors.FromError(err)
		}
	} else {
		err = s.db.Transaction(func(tx *gorm.DB) error {
			bookRepo := s.bookingRepo.WithTx(tx)

			slot, slotErr := bookRepo.FindAndLockSlot(booking.AppCode, booking.Date, booking.StartTime)
			if slotErr != nil {
				return slotErr
			}

			slot.SeatsBooked -= booking.AttendeeCount
			normalizeSlotState(slot)
			if updateErr := bookRepo.Update(slot); updateErr != nil {
				return updateErr
			}

			booking.Available = true
			booking.Status = "rejected"
			if updateErr := bookRepo.Update(booking); updateErr != nil {
				return updateErr
			}

			return nil
		})

		if err != nil {
			log.Printf("[RejectBooking] DB error: %v", err)
			return nil, serviceerrors.FromError(err)
		}
	}

	if err == nil {
		go func() {
			if err := s.notificationService.SendBookingRejection(booking); err != nil {
				s.bookingRepo.UpdateNotificationStatus(booking.ID, "failed", "email")
			} else {
				s.bookingRepo.UpdateNotificationStatus(booking.ID, "sent", "email")
			}
		}()

		message := fmt.Sprintf("Booking %s was rejected.", booking.BookingCode)
		s.eventNotificationService.CreateEventNotification(appointment.OwnerID, "BOOKING_REJECTED", message, booking.ID)
	}

	return booking, nil
}
