package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/db"
	myerrors "github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/m13ha/appointment_master/utils"
	"github.com/morkid/paginate"
)

func ToAppointmentResponse(appointment *entities.Appointment) *responses.AppointmentResponse {
	return &responses.AppointmentResponse{
		ID:              appointment.ID,
		Title:           appointment.Title,
		StartTime:       appointment.StartTime,
		EndTime:         appointment.EndTime,
		StartDate:       appointment.StartDate,
		EndDate:         appointment.EndDate,
		BookingDuration: appointment.BookingDuration,
		Type:            appointment.Type,
		MaxAttendees:    appointment.MaxAttendees,
		AppCode:         appointment.AppCode,
		CreatedAt:       appointment.CreatedAt,
		UpdatedAt:       appointment.UpdatedAt,
		Description:     appointment.Description,
	}
}

func CreateAppointment(req requests.AppointmentRequest, userId uuid.UUID) (*responses.AppointmentResponse, error) {
	// Validate request
	if err := utils.Validate(req); err != nil {
		return nil, myerrors.NewUserError("Invalid appointment data. Please check your input.")
	}

	if req.EndTime.Before(req.StartTime) {
		return nil, myerrors.NewUserError("End time cannot be before start time.")
	}

	if req.EndDate.Before(req.StartDate) {
		return nil, myerrors.NewUserError("End date cannot be before start date.")
	}

	// Validate booking duration fits within time window
	duration := req.EndTime.Sub(req.StartTime)
	if duration.Minutes() < float64(req.BookingDuration) {
		return nil, myerrors.NewUserError("Booking duration exceeds available time window.")
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

	code, err := generateUniqueAppCode()
	if err != nil {
		log.Printf("[CreateAppointment] Internal error: %v", err)
		return nil, fmt.Errorf("internal error")
	}
	appointment.AppCode = code

	if err := db.DB.Create(appointment).Error; err != nil {
		log.Printf("[CreateAppointment] DB error: %v", err)
		return nil, fmt.Errorf("internal error")
	}

	return ToAppointmentResponse(appointment), nil
}

func GetAllAppointmentsCreatedByUser(userID string, r *http.Request) (any, error) {
	query := db.DB.Model(&entities.Appointment{}).Where("owner_id = ?", userID)
	if r == nil {
		var appointments []entities.Appointment
		if err := query.Find(&appointments).Error; err != nil {
			log.Printf("[GetAllAppointmentsCreatedByUser] DB error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
		return appointments, nil
	}
	p := paginate.New()
	result := p.With(query).Request(r).Response(&[]responses.AppointmentResponse{})
	return &result, nil
}

func isAppCodeAvailable(appCode string) (bool, error) {
	var appointment entities.Appointment
	err := db.DB.Where("app_code = ?", appCode).Order("end_date desc").First(&appointment).Error
	if err != nil {
		return true, nil
	}
	now := time.Now()
	holdUntil := appointment.EndDate.Add(14 * 24 * time.Hour)
	if now.Before(holdUntil) {
		return false, nil
	}
	return true, nil
}

func generateUniqueAppCode() (string, error) {
	for i := 0; i < 10; i++ {
		code := utils.GenerateAppCode()
		available, err := isAppCodeAvailable(code)
		if err != nil {
			return "", err
		}
		if available {
			return code, nil
		}
	}
	return "", fmt.Errorf("could not generate unique AppCode after 10 attempts")
}
