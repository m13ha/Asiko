package api

import (
	"github.com/gin-gonic/gin"
	"github.com/m13ha/appointment_master/middleware"
	"github.com/m13ha/appointment_master/services"
)

type Handler struct {
	userService        services.UserService
	appointmentService services.AppointmentService
	bookingService     services.BookingService
}

func NewHandler(userService services.UserService, appointmentService services.AppointmentService, bookingService services.BookingService) *Handler {
	return &Handler{
		userService:        userService,
		appointmentService: appointmentService,
		bookingService:     bookingService,
	}
}

func RegisterRoutes(r *gin.Engine, userService services.UserService, appointmentService services.AppointmentService, bookingService services.BookingService) {
	h := NewHandler(userService, appointmentService, bookingService)

	r.POST("/login", h.Login)
	r.POST("/logout", h.Logout)
	r.POST("/users", h.CreateUser)
	r.POST("/appointments/book", h.BookGuestAppointment)
	r.GET("/appointments/slots/:id", h.GetAvailableSlots)
	r.GET("/appointments/slots/:id/by-day", h.GetAvailableSlotsByDay)
	r.GET("/bookings/:booking_code", h.GetBookingByCodeHandler)
	r.PUT("/bookings/:booking_code", h.UpdateBookingByCodeHandler)
	r.DELETE("/bookings/:booking_code", h.CancelBookingByCodeHandler)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/appointments", h.CreateAppointment)
		protected.GET("/appointments/users/:id", h.GetUsersRegisteredForAppointment)
		protected.GET("/appointments/my", h.GetAppointmentsCreatedByUser)
		protected.GET("/appointments/registered", h.GetUserRegisteredBookings)
		protected.POST("/appointments/book/registered", h.BookRegisteredUserAppointment)
	}
}
