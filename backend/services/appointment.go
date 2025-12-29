package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	serviceerrors "github.com/m13ha/asiko/errors/serviceerrors"
	"github.com/m13ha/asiko/events"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	"github.com/m13ha/asiko/models/responses"
	"github.com/m13ha/asiko/repository"
	"github.com/m13ha/asiko/utils"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type appointmentServiceImpl struct {
	appointmentRepo          repository.AppointmentRepository
	bookingRepo              repository.BookingRepository
	userRepo                 repository.UserRepository
	eventBus                events.EventBus
	eventNotificationService EventNotificationService
	db                       *gorm.DB
}

type StatusRefreshSummary struct {
	PendingToOngoing int64
	Completed        int64
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository, bookingRepo repository.BookingRepository, userRepo repository.UserRepository, eventBus events.EventBus, eventNotificationService EventNotificationService, db *gorm.DB) AppointmentService {
	return &appointmentServiceImpl{
		appointmentRepo:          appointmentRepo,
		bookingRepo:              bookingRepo,
		userRepo:                 userRepo,
		eventBus:                 eventBus,
		eventNotificationService: eventNotificationService,
		db:                       db,
	}
}

func (s *appointmentServiceImpl) CreateAppointment(req requests.AppointmentRequest, userId uuid.UUID) (*entities.Appointment, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	appointment := &entities.Appointment{
		Title:             req.Title,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		BookingDuration:   req.BookingDuration,
		Type:              entities.AppointmentType(utils.NormalizeString(fmt.Sprintf("%v", req.Type))),
		MaxAttendees:      req.MaxAttendees,
		OwnerID:           userId,
		Description:       req.Description,
		AntiScalpingLevel: req.AntiScalpingLevel,
		Status:            entities.AppointmentStatusPending,
	}

	if err := s.appointmentRepo.Create(appointment); err != nil {
		log.Printf("[CreateAppointment] DB error: %v", err)
		return nil, serviceerrors.FromError(err)
	}

	payload := events.AppointmentEventData{
		Appointment:      appointment,
		OwnerID:          appointment.OwnerID,
		AppointmentTitle: appointment.Title,
	}
	if owner, ownerErr := s.userRepo.FindByID(appointment.OwnerID.String()); ownerErr == nil {
		payload.RecipientEmail = owner.Email
		payload.RecipientName = owner.Name
	}
	if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventAppointmentCreated, Data: payload}); pubErr != nil {
		log.Printf("Failed to publish appointment created event: %v", pubErr)
	}

	return appointment, nil
}

func (s *appointmentServiceImpl) UpdateAppointment(ctx context.Context, appointmentID uuid.UUID, ownerID uuid.UUID, req requests.AppointmentRequest) (*entities.Appointment, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	appointment, err := s.appointmentRepo.FindByIDAndOwner(ctx, appointmentID, ownerID)
	if err != nil {
		return nil, err
	}

	hasBookings, err := s.bookingRepo.HasActiveBookings(appointment.ID)
	if err != nil {
		return nil, serviceerrors.FromError(err)
	}
	if hasBookings {
		return nil, serviceerrors.ConflictError("appointments with booked slots cannot be edited")
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		appRepo := s.appointmentRepo.WithTx(tx)
		bookRepo := s.bookingRepo.WithTx(tx)

		appointment.Title = req.Title
		appointment.StartTime = req.StartTime
		appointment.EndTime = req.EndTime
		appointment.StartDate = req.StartDate
		appointment.EndDate = req.EndDate
		appointment.BookingDuration = req.BookingDuration
		appointment.Type = entities.AppointmentType(utils.NormalizeString(string(req.Type)))
		appointment.MaxAttendees = req.MaxAttendees
		appointment.Description = req.Description
		appointment.AntiScalpingLevel = req.AntiScalpingLevel
		appointment.AttendeesBooked = 0

		if err := appRepo.Update(appointment); err != nil {
			return err
		}

		if err := bookRepo.DeleteSlotsByAppointmentID(appointment.ID); err != nil {
			return err
		}

		slots := appointment.GenerateBookings()
		if len(slots) > 0 {
			if err := tx.Create(&slots).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, serviceerrors.FromError(err)
	}

	payload := events.AppointmentEventData{
		Appointment:      appointment,
		OwnerID:          appointment.OwnerID,
		AppointmentTitle: appointment.Title,
	}
	if owner, ownerErr := s.userRepo.FindByID(appointment.OwnerID.String()); ownerErr == nil {
		payload.RecipientEmail = owner.Email
		payload.RecipientName = owner.Name
	}
	if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventAppointmentUpdated, Data: payload}); pubErr != nil {
		log.Printf("Failed to publish appointment updated event: %v", pubErr)
	}

	return appointment, nil
}

func (s *appointmentServiceImpl) DeleteAppointment(ctx context.Context, appointmentID uuid.UUID, ownerID uuid.UUID) (*entities.Appointment, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	appointment, err := s.appointmentRepo.FindByIDAndOwner(ctx, appointmentID, ownerID)
	if err != nil {
		return nil, err
	}

	var affectedBookings []entities.Booking
	var notifyBookings []entities.Booking
	err = s.db.Transaction(func(tx *gorm.DB) error {
		appRepo := s.appointmentRepo.WithTx(tx)
		bookRepo := s.bookingRepo.WithTx(tx)

		bookings, err := bookRepo.GetActiveBookingsForAppointment(appointment.ID)
		if err != nil {
			return err
		}
		affectedBookings = bookings

		lockedAppointment := appointment
		if appointment.Type == entities.Party {
			locked, err := appRepo.FindAndLock(appointment.AppCode, tx)
			if err != nil {
				return err
			}
			lockedAppointment = locked
		}

		for i := range affectedBookings {
			booking := &affectedBookings[i]

			switch {
			case appointment.Type == entities.Party:
				lockedAppointment.AttendeesBooked -= booking.AttendeeCount
				if lockedAppointment.AttendeesBooked < 0 {
					lockedAppointment.AttendeesBooked = 0
				}
				booking.Status = entities.BookingStatusCancelled
				if err := bookRepo.Update(booking); err != nil {
					return err
				}
				notifyBookings = append(notifyBookings, *booking)
			case booking.IsSlot:
				if appointment.Type == entities.Group {
					booking.SeatsBooked = 0
					booking.Status = entities.BookingStatusCancelled
					booking.NormalizeState()
					booking.Available = false
					if err := bookRepo.Update(booking); err != nil {
						return err
					}
					continue
				}
				booking.SeatsBooked = 0
				booking.Status = entities.BookingStatusCancelled
				booking.NormalizeState()
				booking.Available = false
				if err := bookRepo.Update(booking); err != nil {
					return err
				}
				notifyBookings = append(notifyBookings, *booking)
			default:
				slot, err := bookRepo.FindAndLockSlot(booking.AppCode, booking.Date, booking.StartTime)
				if err != nil {
					return err
				}
				slot.SeatsBooked -= booking.AttendeeCount
				slot.NormalizeState()
				if err := bookRepo.Update(slot); err != nil {
					return err
				}
				booking.Available = true
				booking.Status = entities.BookingStatusCancelled
				if err := bookRepo.Update(booking); err != nil {
					return err
				}
				notifyBookings = append(notifyBookings, *booking)
			}
		}

		if appointment.Type == entities.Party {
			if err := appRepo.Update(lockedAppointment); err != nil {
				return err
			}
			appointment = lockedAppointment
		}

		appointment.Status = entities.AppointmentStatusCanceled
		if err := appRepo.Update(appointment); err != nil {
			return err
		}
		if err := tx.Delete(appointment).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, serviceerrors.FromError(err)
	}

	for _, booking := range notifyBookings {
		payload := events.BookingEventData{
			Booking:          &booking,
			OwnerID:          appointment.OwnerID,
			AppointmentTitle: appointment.Title,
			RecipientEmail:   booking.Email,
			RecipientName:    booking.Name,
		}
		if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventBookingCancelled, Data: payload}); pubErr != nil {
			log.Printf("Failed to publish booking cancelled event: %v", pubErr)
		}
	}

	payload := events.AppointmentEventData{
		Appointment:      appointment,
		OwnerID:          appointment.OwnerID,
		AppointmentTitle: appointment.Title,
	}
	if owner, ownerErr := s.userRepo.FindByID(appointment.OwnerID.String()); ownerErr == nil {
		payload.RecipientEmail = owner.Email
		payload.RecipientName = owner.Name
	}
	if pubErr := s.eventBus.Publish(context.Background(), events.Event{Name: events.EventAppointmentDeleted, Data: payload}); pubErr != nil {
		log.Printf("Failed to publish appointment deleted event: %v", pubErr)
	}

	return appointment, nil
}

func (s *appointmentServiceImpl) GetAllAppointmentsCreatedByUser(userID string, r *http.Request, statuses []entities.AppointmentStatus) paginate.Page {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return paginate.New().With(nil).Request(r).Response(&[]responses.AppointmentResponse{})
	}
	var ctx context.Context
	if r != nil {
		ctx = r.Context()
	} else {
		ctx = context.Background()
	}
	return s.appointmentRepo.GetAppointmentsByOwnerIDQuery(ctx, r, uid, statuses)
}

func (s *appointmentServiceImpl) CancelAppointment(ctx context.Context, appointmentID uuid.UUID, ownerID uuid.UUID) (*entities.Appointment, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	appointment, err := s.appointmentRepo.FindByIDAndOwner(ctx, appointmentID, ownerID)
	if err != nil {
		return nil, err
	}

	if appointment.Status == entities.AppointmentStatusCanceled {
		return nil, serviceerrors.ConflictError("appointment already canceled")
	}

	if appointment.Status == entities.AppointmentStatusCompleted {
		return nil, serviceerrors.ConflictError("finished appointments cannot be canceled")
	}

	if !entities.CanTransitionAppointmentStatus(appointment.Status, entities.AppointmentStatusCanceled) {
		return nil, serviceerrors.ConflictError("appointment cannot be canceled in its current status")
	}

	if err := s.appointmentRepo.UpdateStatus(ctx, appointment.ID, entities.AppointmentStatusCanceled); err != nil {
		return nil, err
	}
	appointment.Status = entities.AppointmentStatusCanceled

	message := fmt.Sprintf("Appointment '%s' was canceled.", appointment.Title)
	s.eventNotificationService.CreateEventNotification(ownerID, "APPOINTMENT_CANCELED", message, appointment.ID)

	return appointment, nil
}

func (s *appointmentServiceImpl) RefreshStatuses(ctx context.Context, now time.Time) (StatusRefreshSummary, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if now.IsZero() {
		now = time.Now()
	}

	var summary StatusRefreshSummary
	updated, err := s.appointmentRepo.MarkAppointmentsOngoing(ctx, now)
	if err != nil {
		return summary, err
	}
	summary.PendingToOngoing = updated

	updated, err = s.appointmentRepo.MarkAppointmentsCompleted(ctx, now)
	if err != nil {
		return summary, err
	}
	summary.Completed = updated

	return summary, nil
}

func (s *appointmentServiceImpl) GetAppointmentByAppCode(appCode string) (*entities.Appointment, error) {
	appointment, err := s.appointmentRepo.FindAppointmentByAppCode(appCode)
	if err != nil {
		return nil, err
	}
	return appointment, nil
}
