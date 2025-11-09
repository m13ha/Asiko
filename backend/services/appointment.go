package services

import (
    "fmt"
    "context"
    "log"
    "net/http"

	"github.com/google/uuid"
    myerrors "github.com/m13ha/asiko/errors"
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
	}

    if err := s.appointmentRepo.Create(appointment); err != nil {
        log.Printf("[CreateAppointment] DB error: %v", err)
        return nil, myerrors.FromError(err)
    }

	message := fmt.Sprintf("New appointment '%s' created.", appointment.Title)
	s.eventNotificationService.CreateEventNotification(appointment.OwnerID, "APPOINTMENT_CREATED", message, appointment.ID)

	return appointment, nil
}

func (s *appointmentServiceImpl) GetAllAppointmentsCreatedByUser(userID string, r *http.Request) paginate.Page {
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
    return s.appointmentRepo.GetAppointmentsByOwnerIDQuery(ctx, uid)
}
