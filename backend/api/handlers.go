package api

import (
	"github.com/gin-gonic/gin"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/services"
)

type Handler struct {
	userService              services.UserService
	appointmentService       services.AppointmentService
	bookingService           services.BookingService
	analyticsService         services.AnalyticsService
	banService               services.BanListService
	eventNotificationService services.EventNotificationService
}

func NewHandler(userService services.UserService, appointmentService services.AppointmentService, bookingService services.BookingService, analyticsService services.AnalyticsService, banServices services.BanListService, eventNotificationService services.EventNotificationService) *Handler {
	return &Handler{
		userService:              userService,
		appointmentService:       appointmentService,
		bookingService:           bookingService,
		analyticsService:         analyticsService,
		banService:               banServices,
		eventNotificationService: eventNotificationService,
	}
}

func RegisterRoutes(r *gin.Engine, userService services.UserService, appointmentService services.AppointmentService, bookingService services.BookingService, analyticsService services.AnalyticsService, banServices services.BanListService, eventNotificationService services.EventNotificationService) {
	h := NewHandler(userService, appointmentService, bookingService, analyticsService, banServices, eventNotificationService)

	r.POST("/login", h.Login)
	r.POST("/logout", h.Logout)
	r.POST("/users", h.CreateUser)
	r.POST("/auth/verify-registration", h.VerifyRegistrationHandler)
	r.POST("/auth/resend-verification", h.ResendVerificationHandler)
	r.POST("/auth/device-token", h.GenerateDeviceTokenHandler)
	r.POST("/auth/refresh", h.Refresh)
	r.POST("/auth/forgot-password", h.ForgotPasswordHandler)
	r.POST("/auth/reset-password", h.ResetPasswordHandler)
	r.POST("/auth/change-password", middleware.AuthMiddleware(), h.ChangePasswordHandler)
	r.POST("/appointments/book", h.BookGuestAppointment)
	r.GET("/appointments/code/:app_code", h.GetAppointmentByAppCode)
	r.GET("/appointments/slots/:app_code", h.GetAvailableSlots)
	r.GET("/appointments/dates/:app_code", h.GetAvailableDates)
	r.GET("/appointments/slots/:app_code/by-day", h.GetAvailableSlotsByDay)
	r.GET("/bookings/:booking_code", h.GetBookingByCodeHandler)
	r.PUT("/bookings/:booking_code", h.UpdateBookingByCodeHandler)
	r.DELETE("/bookings/:booking_code", h.CancelBookingByCodeHandler)
	r.POST("/bookings/:booking_code/reject", middleware.AuthMiddleware(), h.RejectBookingHandler)

	// Protected routes with authentication middleware
	r.POST("/appointments", middleware.AuthMiddleware(), h.CreateAppointment)
	r.GET("/appointments/my", middleware.AuthMiddleware(), h.GetAppointmentsCreatedByUser)
	r.GET("/appointments/registered", middleware.AuthMiddleware(), h.GetUserRegisteredBookings)
	r.GET("/appointments/users/:app_code", middleware.AuthMiddleware(), h.GetUsersRegisteredForAppointment)
	r.POST("/appointments/book/registered", middleware.AuthMiddleware(), h.BookRegisteredUserAppointment)
	r.GET("/analytics", middleware.AuthMiddleware(), h.GetUserAnalytics)

	banList := r.Group("/ban-list", middleware.AuthMiddleware())
	{
		banList.POST("", h.AddToBanList)
		banList.DELETE("", h.RemoveFromBanList)
		banList.GET("", h.GetBanList)
	}

	notifications := r.Group("/notifications", middleware.AuthMiddleware())
	{
		notifications.GET("", h.GetNotificationsHandler)
		notifications.GET("/unread-count", h.GetUnreadNotificationsCountHandler)
		notifications.PUT("/read-all", h.MarkAllNotificationsAsReadHandler)
	}
}
