package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/m13ha/asiko/errors"
	serviceerrors "github.com/m13ha/asiko/errors/serviceerrors"
	"github.com/m13ha/asiko/events"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	"github.com/m13ha/asiko/repository"
	"github.com/m13ha/asiko/utils"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type bookingServiceImpl struct {
	bookingRepo     repository.BookingRepository
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository
	banListRepo     repository.BanListRepository
	eventBus        events.EventBus
	db              *gorm.DB
}

type BookingStatusRefreshSummary struct {
	Ongoing int64
	Expired int64
}

func NewBookingService(bookingRepo repository.BookingRepository, appointmentRepo repository.AppointmentRepository, userRepo repository.UserRepository, banListRepo repository.BanListRepository, eventBus events.EventBus, db *gorm.DB) BookingService {
	return &bookingServiceImpl{bookingRepo: bookingRepo, appointmentRepo: appointmentRepo, userRepo: userRepo, banListRepo: banListRepo, eventBus: eventBus, db: db}
}

func isRepoNotFound(err error) bool {
	appErr := appErrors.FromAppError(err)
	return appErr != nil && appErr.Code == appErrors.CodeRepoNotFoundError
}

func sameDay(a time.Time, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

func sameClock(a time.Time, b time.Time) bool {
	return a.Hour() == b.Hour() && a.Minute() == b.Minute() && a.Second() == b.Second() && a.Nanosecond() == b.Nanosecond()
}

// performAntiScalpingChecks runs validation based on the appointment's settings.
// It returns the trusted device ID (if applicable) or an error if a check fails.
func (s *bookingServiceImpl) performAntiScalpingChecks(appointment *entities.Appointment, req requests.BookingRequest, bookingEmail string, bookingPhone string) (string, error) {
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
		} else if !isRepoNotFound(err) {
			return "", serviceerrors.FromError(err)
		}
	}

	// Standard check (Email) - runs for both 'standard' and 'strict'
	if level == entities.ScalpingStandard || level == entities.ScalpingStrict {
		normalizedEmail := utils.NormalizeEmail(bookingEmail)
		if normalizedEmail != "" {
			if _, err := s.bookingRepo.FindActiveBookingByEmail(appointment.ID, normalizedEmail); err == nil {
				return "", serviceerrors.ConflictError("this email has already been used to book for this appointment")
			} else if !isRepoNotFound(err) {
				return "", serviceerrors.FromError(err)
			}
		} else if strings.TrimSpace(bookingPhone) != "" {
			normalizedPhone := strings.TrimSpace(bookingPhone)
			if _, err := s.bookingRepo.FindActiveBookingByPhone(appointment.ID, normalizedPhone); err == nil {
				return "", serviceerrors.ConflictError("this phone has already been used to book for this appointment")
			} else if !isRepoNotFound(err) {
				return "", serviceerrors.FromError(err)
			}
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
	if ok, reason := appointment.IsBookable(); !ok {
		return nil, serviceerrors.ConflictError(reason)
	}

	if appointment.Type == entities.Single && req.AttendeeCount != 1 {
		return nil, serviceerrors.ValidationError("single appointments allow exactly one attendee per slot")
	}

	if appointment.Type == entities.Party {
		if !sameDay(req.Date, appointment.StartDate) {
			return nil, serviceerrors.ValidationError("party bookings must use the appointment start date")
		}
		if !sameClock(req.StartTime, appointment.StartTime) || !sameClock(req.EndTime, appointment.EndTime) {
			return nil, serviceerrors.ValidationError("party booking times must match the appointment window")
		}
	}

	// --- 3. Get Booker's Info ---
	var user *entities.User
	var bookingEmail string
	var bookingPhone string
	if userIDStr != "" {
		user, err = s.userRepo.FindByID(userIDStr)
		if err != nil {
			return nil, err
		}
		bookingEmail = user.Email
		if user.PhoneNumber != nil {
			bookingPhone = *user.PhoneNumber
		}
	} else {
		bookingEmail = utils.NormalizeEmail(req.Email)
		bookingPhone = req.Phone
	}

	// --- 4. Ban List Check ---
	normalizedEmail := utils.NormalizeEmail(bookingEmail)
	if normalizedEmail != "" {
		if _, banErr := s.banListRepo.FindByUserAndEmail(appointment.OwnerID, normalizedEmail); banErr == nil {
			return nil, serviceerrors.ForbiddenError("you are not allowed to book this appointment")
		} else if !isRepoNotFound(banErr) {
			return nil, serviceerrors.FromError(banErr)
		}
	}

	// --- 4. Anti-Scalping Checks ---
	trustedDeviceID, err := s.performAntiScalpingChecks(appointment, req, bookingEmail, bookingPhone)
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
	status := entities.BookingStatusPending
	err := s.db.Transaction(func(tx *gorm.DB) error {
		appRepo := s.appointmentRepo.WithTx(tx)
		bookRepo := s.bookingRepo.WithTx(tx)

		lockedAppointment, err := appRepo.FindAndLock(req.AppCode, tx)
		if err != nil {
			return err
		}
		if ok, reason := lockedAppointment.IsBookable(); !ok {
			return serviceerrors.ConflictError(reason)
		}

		if lockedAppointment.AttendeesBooked+req.AttendeeCount > lockedAppointment.MaxAttendees {
			return serviceerrors.BookingCapacityExceededError("not enough capacity for this party")
		}

		startDateTime := time.Date(
			lockedAppointment.StartDate.Year(), lockedAppointment.StartDate.Month(), lockedAppointment.StartDate.Day(),
			lockedAppointment.StartTime.Hour(), lockedAppointment.StartTime.Minute(), 0, 0,
			time.UTC,
		)
		endDateTime := time.Date(
			lockedAppointment.EndDate.Year(), lockedAppointment.EndDate.Month(), lockedAppointment.EndDate.Day(),
			lockedAppointment.EndTime.Hour(), lockedAppointment.EndTime.Minute(), 0, 0,
			time.UTC,
		)

		booking = &entities.Booking{
			AppointmentID: lockedAppointment.ID,
			AppCode:       lockedAppointment.AppCode,
			Date:          lockedAppointment.StartDate.UTC(),
			StartTime:     startDateTime,
			EndTime:       endDateTime,
			Available:     false,
			IsSlot:        false,
			Capacity:      req.AttendeeCount,
			SeatsBooked:   req.AttendeeCount,
			AttendeeCount: req.AttendeeCount,
			Description:   req.Description,
			Status:        status,
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
		// Publish event
		payload := events.BookingEventData{
			Booking:          booking,
			OwnerID:          appointment.OwnerID,
			AppointmentTitle: appointment.Title,
			RecipientEmail:   booking.Email,
			RecipientName:    booking.Name,
		}
		if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventBookingCreated, Data: payload}); pubErr != nil {
			log.Printf("Failed to publish booking created event: %v", pubErr)
		}
	}

	return booking, err
}

func (s *bookingServiceImpl) bookSlotAppointment(req requests.BookingRequest, user *entities.User, appointment *entities.Appointment, deviceID string) (*entities.Booking, error) {
	status := entities.BookingStatusPending
	if appointment.Type == entities.Group {
		var reservation *entities.Booking
		err := s.db.Transaction(func(tx *gorm.DB) error {
			appRepo := s.appointmentRepo.WithTx(tx)
			bookRepo := s.bookingRepo.WithTx(tx)

			lockedAppointment, err := appRepo.FindAndLock(req.AppCode, tx)
			if err != nil {
				return err
			}
			if ok, reason := lockedAppointment.IsBookable(); !ok {
				return serviceerrors.ConflictError(reason)
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
				Status:        status,
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
			lockedSlot.NormalizeState()
			if lockedSlot.SeatsBooked > 0 {
				lockedSlot.Status = entities.BookingStatusPending
			} else {
				lockedSlot.Status = entities.BookingStatusActive
			}

			if err := bookRepo.Update(lockedSlot); err != nil {
				log.Printf("[bookSlot] failed to update slot: %v", err)
				return serviceerrors.FromError(err)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		// Publish event
		payload := events.BookingEventData{
			Booking:          reservation,
			OwnerID:          appointment.OwnerID,
			AppointmentTitle: appointment.Title,
			RecipientEmail:   reservation.Email,
			RecipientName:    reservation.Name,
		}
		if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventBookingCreated, Data: payload}); pubErr != nil {
			log.Printf("Failed to publish booking created event: %v", pubErr)
		}

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
		if ok, reason := lockedAppointment.IsBookable(); !ok {
			return serviceerrors.ConflictError(reason)
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
		lockedSlot.Status = status
		lockedSlot.NormalizeState()

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

	if err == nil {
		// Publish event
		payload := events.BookingEventData{
			Booking:          slot,
			OwnerID:          appointment.OwnerID,
			AppointmentTitle: appointment.Title,
			RecipientEmail:   slot.Email,
			RecipientName:    slot.Name,
		}
		if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventBookingCreated, Data: payload}); pubErr != nil {
			log.Printf("Failed to publish booking created event: %v", pubErr)
		}
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
	if strings.ToLower(booking.Status) == entities.BookingStatusOngoing {
		return nil, serviceerrors.ConflictError("ongoing bookings cannot be rescheduled")
	}

	appointment, appErr := s.appointmentRepo.FindAppointmentByAppCode(booking.AppCode)
	if appErr != nil {
		return nil, serviceerrors.FromError(appErr)
	}
	if appointment.Type == entities.Party {
		if req.AppCode != booking.AppCode {
			return nil, serviceerrors.ValidationError("party bookings cannot be moved to another appointment")
		}
		if !sameDay(req.Date, appointment.StartDate) || !sameClock(req.StartTime, appointment.StartTime) || !sameClock(req.EndTime, appointment.EndTime) {
			return nil, serviceerrors.ValidationError("party booking times must match the appointment window")
		}
		if req.AttendeeCount < 1 {
			return nil, serviceerrors.ValidationError("attendee count must be at least 1")
		}
		err = s.db.Transaction(func(tx *gorm.DB) error {
			appRepo := s.appointmentRepo.WithTx(tx)
			bookRepo := s.bookingRepo.WithTx(tx)

			lockedAppointment, err := appRepo.FindAndLock(appointment.AppCode, tx)
			if err != nil {
				return err
			}
			delta := req.AttendeeCount - booking.AttendeeCount
			if lockedAppointment.AttendeesBooked+delta > lockedAppointment.MaxAttendees {
				return serviceerrors.BookingCapacityExceededError("not enough capacity for this party")
			}
			lockedAppointment.AttendeesBooked += delta
			if err := appRepo.Update(lockedAppointment); err != nil {
				return err
			}

			booking.AttendeeCount = req.AttendeeCount
			booking.Description = req.Description
			if updateErr := bookRepo.Update(booking); updateErr != nil {
				return updateErr
			}

			return nil
		})
		if err != nil {
			return nil, serviceerrors.FromError(err)
		}

		payload := events.BookingEventData{
			Booking:          booking,
			OwnerID:          appointment.OwnerID,
			AppointmentTitle: appointment.Title,
			RecipientEmail:   booking.Email,
			RecipientName:    booking.Name,
		}
		if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventBookingUpdated, Data: payload}); pubErr != nil {
			log.Printf("Failed to publish booking updated event: %v", pubErr)
		}
		return booking, nil
	}

	if booking.IsSlot && req.AttendeeCount != 1 {
		return nil, serviceerrors.ValidationError("single-slot bookings only support one attendee")
	}

	currentCount := booking.AttendeeCount
	if booking.IsSlot && currentCount < 1 {
		currentCount = booking.Capacity
	}

	wasConfirmed := strings.ToLower(booking.Status) == entities.BookingStatusConfirmed
	err = s.db.Transaction(func(tx *gorm.DB) error {
		bookRepo := s.bookingRepo.WithTx(tx)
		capacityErr := serviceerrors.BookingCapacityExceededError("not enough capacity for this slot")
		sameSlot := booking.AppCode == req.AppCode && booking.Date.Equal(req.Date) && booking.StartTime.Equal(req.StartTime)

		if sameSlot {
			slot, slotErr := bookRepo.FindAndLockSlot(req.AppCode, req.Date, req.StartTime)
			if slotErr != nil {
				return slotErr
			}
			delta := req.AttendeeCount - currentCount
			if delta > 0 && slot.SeatsBooked+delta > slot.Capacity {
				return capacityErr
			}
			slot.SeatsBooked += delta
			slot.NormalizeState()
			if !booking.IsSlot {
				if slot.SeatsBooked > 0 {
					slot.Status = entities.BookingStatusPending
				} else {
					slot.Status = entities.BookingStatusActive
				}
			}
			if strings.ToLower(booking.Status) == entities.BookingStatusPending {
				slot.Status = booking.Status
			}
			if updateErr := bookRepo.Update(slot); updateErr != nil {
				return updateErr
			}
		} else {
			oldSlot, slotErr := bookRepo.FindAndLockSlot(booking.AppCode, booking.Date, booking.StartTime)
			if slotErr != nil {
				return slotErr
			}

			newSlot, slotErr := bookRepo.FindAndLockSlot(req.AppCode, req.Date, req.StartTime)
			if slotErr != nil {
				return slotErr
			}

			if booking.IsSlot {
				if newSlot.SeatsBooked+req.AttendeeCount > newSlot.Capacity {
					return capacityErr
				}

				// Swap booking ownership to the new slot and release the old slot.
				oldCode := oldSlot.BookingCode
				newCode := newSlot.BookingCode
				tempCode := "TMP-" + uuid.NewString()

				// Move the old slot off its booking code first to avoid unique conflicts.
				oldSlot.BookingCode = tempCode
				oldSlot.UserID = nil
				oldSlot.Name = ""
				oldSlot.Email = ""
				oldSlot.Phone = ""
				oldSlot.Description = ""
				oldSlot.DeviceID = ""
				oldSlot.Status = entities.BookingStatusActive
				oldSlot.SeatsBooked = 0
				oldSlot.NormalizeState()
				if updateErr := bookRepo.Update(oldSlot); updateErr != nil {
					return updateErr
				}

				newSlot.BookingCode = oldCode
				newSlot.UserID = booking.UserID
				newSlot.Name = booking.Name
				newSlot.Email = booking.Email
				newSlot.Phone = booking.Phone
				newSlot.Description = req.Description
				newSlot.DeviceID = booking.DeviceID
				newSlot.SeatsBooked = req.AttendeeCount
				newSlot.AttendeeCount = req.AttendeeCount
				newSlot.Available = false
				if wasConfirmed {
					newSlot.Status = entities.BookingStatusPending
				} else if strings.ToLower(booking.Status) == entities.BookingStatusPending {
					newSlot.Status = booking.Status
				} else {
					newSlot.Status = entities.BookingStatusActive
				}
				if updateErr := bookRepo.Update(newSlot); updateErr != nil {
					return updateErr
				}

				oldSlot.BookingCode = newCode
				if updateErr := bookRepo.Update(oldSlot); updateErr != nil {
					return updateErr
				}

				booking = newSlot
			} else {
				oldSlot.SeatsBooked -= booking.AttendeeCount
				oldSlot.NormalizeState()
				if oldSlot.SeatsBooked > 0 {
					oldSlot.Status = entities.BookingStatusPending
				} else {
					oldSlot.Status = entities.BookingStatusActive
				}
				if updateErr := bookRepo.Update(oldSlot); updateErr != nil {
					return updateErr
				}

				if newSlot.SeatsBooked+req.AttendeeCount > newSlot.Capacity {
					return capacityErr
				}
				newSlot.SeatsBooked += req.AttendeeCount
				newSlot.NormalizeState()
				if newSlot.SeatsBooked > 0 {
					newSlot.Status = entities.BookingStatusPending
				} else {
					newSlot.Status = entities.BookingStatusActive
				}
				if updateErr := bookRepo.Update(newSlot); updateErr != nil {
					return updateErr
				}
			}
		}

		if !booking.IsSlot {
			booking.AppCode = req.AppCode
			booking.Date = req.Date
			booking.StartTime = req.StartTime
			booking.EndTime = req.EndTime
			booking.AttendeeCount = req.AttendeeCount
			booking.Description = req.Description
			booking.Available = false
			if !sameSlot && wasConfirmed {
				booking.Status = entities.BookingStatusPending
			}

			if updateErr := bookRepo.Update(booking); updateErr != nil {
				return updateErr
			}
		}

		return nil
	})

	if err != nil {
		return nil, serviceerrors.FromError(err)
	}

	// Publish update event
	payload := events.BookingEventData{
		Booking:          booking,
		OwnerID:          appointment.OwnerID,
		AppointmentTitle: appointment.Title,
		RecipientEmail:   booking.Email,
		RecipientName:    booking.Name,
	}
	if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventBookingUpdated, Data: payload}); pubErr != nil {
		log.Printf("Failed to publish booking updated event: %v", pubErr)
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

	if !entities.CanTransitionBookingStatus(booking.Status, entities.BookingStatusCancelled) {
		return nil, serviceerrors.ConflictError("booking cannot be cancelled in its current status")
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

			booking.Status = entities.BookingStatusCancelled
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

			decrement := booking.AttendeeCount
			if booking.IsSlot && decrement < 1 {
				decrement = booking.Capacity
			}
			slot.SeatsBooked -= decrement
			slot.NormalizeState()
			if appointment.Type == entities.Group {
				if slot.SeatsBooked > 0 {
					slot.Status = entities.BookingStatusPending
				} else {
					slot.Status = entities.BookingStatusActive
				}
			}
			if updateErr := bookRepo.Update(slot); updateErr != nil {
				return updateErr
			}

			booking.Available = true
			booking.Status = entities.BookingStatusCancelled
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
		// Publish cancellation event
		payload := events.BookingEventData{
			Booking:          booking,
			OwnerID:          appointment.OwnerID,
			AppointmentTitle: appointment.Title,
			RecipientEmail:   booking.Email,
			RecipientName:    booking.Name,
		}
		if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventBookingCancelled, Data: payload}); pubErr != nil {
			log.Printf("Failed to publish booking cancelled event: %v", pubErr)
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
		return nil, serviceerrors.FromError(err)
	}

	if appointment.OwnerID != ownerID {
		return nil, serviceerrors.ForbiddenError("you are not the owner of this appointment")
	}

	if !entities.CanTransitionBookingStatus(booking.Status, entities.BookingStatusRejected) {
		return nil, serviceerrors.ConflictError("booking cannot be rejected in its current status")
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

			booking.Status = entities.BookingStatusRejected
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

			decrement := booking.AttendeeCount
			if booking.IsSlot && decrement < 1 {
				decrement = booking.Capacity
			}
			slot.SeatsBooked -= decrement
			slot.NormalizeState()
			if appointment.Type == entities.Group {
				if slot.SeatsBooked > 0 {
					slot.Status = entities.BookingStatusPending
				} else {
					slot.Status = entities.BookingStatusActive
				}
			}
			if updateErr := bookRepo.Update(slot); updateErr != nil {
				return updateErr
			}

			if booking.IsSlot {
				booking.UserID = nil
				booking.Name = ""
				booking.Email = ""
				booking.Phone = ""
				booking.Description = ""
				booking.DeviceID = ""
				booking.SeatsBooked = 0
				booking.Status = entities.BookingStatusActive
				booking.NormalizeState()
			} else {
				booking.Available = true
				booking.Status = entities.BookingStatusRejected
			}
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
		// Publish rejection event
		payload := events.BookingEventData{
			Booking:          booking,
			OwnerID:          appointment.OwnerID,
			AppointmentTitle: appointment.Title,
			RecipientEmail:   booking.Email,
			RecipientName:    booking.Name,
		}
		if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventBookingRejected, Data: payload}); pubErr != nil {
			log.Printf("Failed to publish booking rejected event: %v", pubErr)
		}
	}

	return booking, nil
}

func (s *bookingServiceImpl) ConfirmBooking(bookingCode string, ownerID uuid.UUID) (*entities.Booking, error) {
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

	if !entities.CanTransitionBookingStatus(booking.Status, entities.BookingStatusConfirmed) {
		return nil, serviceerrors.ConflictError("booking cannot be confirmed in its current status")
	}

	booking.Available = false
	booking.Status = entities.BookingStatusConfirmed
	if err := s.bookingRepo.Update(booking); err != nil {
		return nil, serviceerrors.FromError(err)
	}

	payload := events.BookingEventData{
		Booking:          booking,
		OwnerID:          appointment.OwnerID,
		AppointmentTitle: appointment.Title,
		RecipientEmail:   booking.Email,
		RecipientName:    booking.Name,
	}
	if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventBookingConfirmed, Data: payload}); pubErr != nil {
		log.Printf("Failed to publish booking confirmed event: %v", pubErr)
	}

	return booking, nil
}
