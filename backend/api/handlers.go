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
	analyticsService   services.AnalyticsService
	banService         services.BanListService
}

func NewHandler(userService services.UserService, appointmentService services.AppointmentService, bookingService services.BookingService, analyticsService services.AnalyticsService, banServices services.BanListService) *Handler {
	return &Handler{
		userService:        userService,
		appointmentService: appointmentService,
		bookingService:     bookingService,
		analyticsService:   analyticsService,
		banService:         banServices,
	}
}

func RegisterRoutes(r *gin.Engine, userService services.UserService, appointmentService services.AppointmentService, bookingService services.BookingService, analyticsService services.AnalyticsService, banServices services.BanListService) {
	h := NewHandler(userService, appointmentService, bookingService, analyticsService, banServices)

	r.POST("/login", h.Login)
	r.POST("/logout", h.Logout)
	r.POST("/users", h.CreateUser)
	r.POST("/auth/device-token", h.GenerateDeviceTokenHandler)
	r.POST("/appointments/book", h.BookGuestAppointment)
	r.GET("/appointments/slots/:id", h.GetAvailableSlots)
	r.GET("/appointments/slots/:id/by-day", h.GetAvailableSlotsByDay)
	r.GET("/bookings/:booking_code", h.GetBookingByCodeHandler)
	r.PUT("/bookings/:booking_code", h.UpdateBookingByCodeHandler)
	r.DELETE("/bookings/:booking_code", h.CancelBookingByCodeHandler)
	r.POST("/bookings/:booking_code/reject", middleware.AuthMiddleware(), h.RejectBookingHandler)

	// Protected routes with authentication middleware
	r.POST("/appointments", middleware.AuthMiddleware(), h.CreateAppointment)
	r.GET("/appointments/my", middleware.AuthMiddleware(), h.GetAppointmentsCreatedByUser)
	r.GET("/appointments/registered", middleware.AuthMiddleware(), h.GetUserRegisteredBookings)
	r.GET("/appointments/users/:id", middleware.AuthMiddleware(), h.GetUsersRegisteredForAppointment)
	r.POST("/appointments/book/registered", middleware.AuthMiddleware(), h.BookRegisteredUserAppointment)
	r.GET("/analytics", middleware.AuthMiddleware(), h.GetUserAnalytics)

	banList := r.Group("/ban-list", middleware.AuthMiddleware())
	{
		banList.POST("", h.AddToBanList)
		banList.DELETE("", h.RemoveFromBanList)
		banList.GET("", h.GetBanList)
	}
}
