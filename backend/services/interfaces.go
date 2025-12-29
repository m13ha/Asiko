package services

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	"github.com/m13ha/asiko/models/responses"
	"github.com/morkid/paginate"
)

type UserService interface {
	CreateUser(userReq requests.UserRequest) (*responses.UserResponse, error)
	AuthenticateUser(email, password string) (*entities.User, error)
	VerifyRegistration(email, code string) (string, error)
	ResendVerificationCode(email string) error
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
	ChangePassword(userID, oldPassword, newPassword string) error
}

type BookingService interface {
	BookAppointment(req requests.BookingRequest, userIDStr string) (*entities.Booking, error)
	BookRegisteredUserAppointment(req requests.BookingRequest, userIDStr string) (*entities.Booking, error)
	BookGuestAppointment(req requests.BookingRequest) (*entities.Booking, error)
	GetAllBookingsForAppointment(ctx context.Context, req *http.Request, appcode string) (paginate.Page, error)
	GetUserBookings(ctx context.Context, req *http.Request, userID string, statuses []string) (paginate.Page, error)
	GetAvailableSlots(req *http.Request, appcode string) (paginate.Page, error)
	GetAvailableSlotsByDay(req *http.Request, appcode string, dateStr string) (paginate.Page, error)
	GetAvailableDates(ctx context.Context, appcode string) ([]string, error)
	GetBookingByCode(bookingCode string) (*entities.Booking, error)
	UpdateBookingByCode(bookingCode string, req requests.BookingRequest) (*entities.Booking, error)
	CancelBookingByCode(bookingCode string) (*entities.Booking, error)
	ConfirmBooking(bookingCode string, ownerID uuid.UUID) (*entities.Booking, error)
	RejectBooking(bookingCode string, ownerID uuid.UUID) (*entities.Booking, error)
	RefreshBookingStatuses(ctx context.Context, now time.Time) (BookingStatusRefreshSummary, error)
}

type AppointmentService interface {
	CreateAppointment(req requests.AppointmentRequest, userId uuid.UUID) (*entities.Appointment, error)
	GetAllAppointmentsCreatedByUser(userID string, r *http.Request, statuses []entities.AppointmentStatus) paginate.Page
	UpdateAppointment(ctx context.Context, appointmentID uuid.UUID, ownerID uuid.UUID, req requests.AppointmentRequest) (*entities.Appointment, error)
	DeleteAppointment(ctx context.Context, appointmentID uuid.UUID, ownerID uuid.UUID) (*entities.Appointment, error)
	CancelAppointment(ctx context.Context, appointmentID uuid.UUID, ownerID uuid.UUID) (*entities.Appointment, error)
	RefreshStatuses(ctx context.Context, now time.Time) (StatusRefreshSummary, error)
	GetAppointmentByAppCode(appCode string) (*entities.Appointment, error)
}
