package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/m13ha/appointment_master/repository"
	"github.com/m13ha/appointment_master/utils"
	"github.com/morkid/paginate"
)

type appointmentServiceImpl struct {
	appointmentRepo repository.AppointmentRepository
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository) AppointmentService {
	return &appointmentServiceImpl{appointmentRepo: appointmentRepo}
}

func (s *appointmentServiceImpl) CreateAppointment(req requests.AppointmentRequest, userId uuid.UUID) (*entities.Appointment, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	appointment := &entities.Appointment{
		Title:           req.Title,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		BookingDuration: req.BookingDuration,
		Type:            entities.AppointmentType(utils.NormalizeString(fmt.Sprintf("%v", req.Type))),
		MaxAttendees:    req.MaxAttendees,
		OwnerID:         userId,
		Description:     req.Description,
	}

	if err := s.appointmentRepo.Create(appointment); err != nil {
		log.Printf("[CreateAppointment] DB error: %v", err)
		return nil, fmt.Errorf("internal error")
	}

	return appointment, nil
}

func (s *appointmentServiceImpl) GetAllAppointmentsCreatedByUser(userID string, r *http.Request) paginate.Page {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return paginate.New().With(nil).Request(r).Response(&[]responses.AppointmentResponse{})
	}
	return s.appointmentRepo.GetAppointmentsByOwnerIDQuery(r.Context(), uid)
}
