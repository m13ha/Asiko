package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	serviceerrors "github.com/m13ha/asiko/errors/serviceerrors"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	"github.com/m13ha/asiko/models/responses"
	"github.com/m13ha/asiko/repository"
	"github.com/m13ha/asiko/utils"
	"github.com/morkid/paginate"
)

type appointmentServiceImpl struct {
	appointmentRepo          repository.AppointmentRepository
	eventNotificationService EventNotificationService
}

type StatusRefreshSummary struct {
	PendingToOngoing int64
	Completed        int64
	Expired          int64
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository, eventNotificationService EventNotificationService) AppointmentService {
	return &appointmentServiceImpl{appointmentRepo: appointmentRepo, eventNotificationService: eventNotificationService}
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

	message := fmt.Sprintf("New appointment '%s' created.", appointment.Title)
	s.eventNotificationService.CreateEventNotification(appointment.OwnerID, "APPOINTMENT_CREATED", message, appointment.ID)

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

	if appointment.Status == entities.AppointmentStatusCompleted || appointment.Status == entities.AppointmentStatusExpired {
		return nil, serviceerrors.ConflictError("finished appointments cannot be canceled")
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

	updated, err = s.appointmentRepo.MarkAppointmentsExpired(ctx, now)
	if err != nil {
		return summary, err
	}
	summary.Expired = updated

	return summary, nil
}

func (s *appointmentServiceImpl) GetAppointmentByAppCode(appCode string) (*entities.Appointment, error) {
	appointment, err := s.appointmentRepo.FindAppointmentByAppCode(appCode)
	if err != nil {
		return nil, err
	}
	return appointment, nil
}
