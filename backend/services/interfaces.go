package services

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/morkid/paginate"
)

type UserService interface {
	CreateUser(userReq requests.UserRequest) (*responses.UserResponse, error)
	AuthenticateUser(email, password string) (*entities.User, error)
}

type BookingService interface {
	BookAppointment(req requests.BookingRequest, userIDStr string) (*entities.Booking, error)
	BookRegisteredUserAppointment(req requests.BookingRequest, userIDStr string) (*entities.Booking, error)
	BookGuestAppointment(req requests.BookingRequest) (*entities.Booking, error)

	GetAllBookingsForAppointment(ctx context.Context, appcode string) (paginate.Page, error)
	GetUserBookings(ctx context.Context, userID string) (paginate.Page, error)
	GetAvailableSlots(ctx context.Context, appcode string) (paginate.Page, error)
	GetAvailableSlotsByDay(ctx context.Context, appcode string, dateStr string) (paginate.Page, error)
	GetBookingByCode(bookingCode string) (*entities.Booking, error)
	UpdateBookingByCode(bookingCode string, req requests.BookingRequest) (*entities.Booking, error)
	CancelBookingByCode(bookingCode string) (*entities.Booking, error)
}

type AppointmentService interface {
	CreateAppointment(req requests.AppointmentRequest, userId uuid.UUID) (*entities.Appointment, error)
	GetAllAppointmentsCreatedByUser(userID string, r *http.Request) paginate.Page
}
