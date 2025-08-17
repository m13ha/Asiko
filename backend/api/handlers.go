package api

import (
	"github.com/gin-gonic/gin"
	"github.com/m13ha/appointment_master/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/login", Login)
	r.POST("/logout", Logout)
	r.POST("/users", CreateUser)
	r.POST("/appointments/book", BookGuestAppointment)
	r.GET("/appointments/slots/:id", GetAvailableSlots)
	r.GET("/bookings/:booking_code", GetBookingByCodeHandler)
	r.PUT("/bookings/:booking_code", UpdateBookingByCodeHandler)
	r.DELETE("/bookings/:booking_code", CancelBookingByCodeHandler)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/appointments", CreateAppointment)
		protected.GET("/appointments/users/:id", GetUsersRegisteredForAppointment)
		protected.GET("/appointments/my", GetAppointmentsCreatedByUser)
		protected.GET("/appointments/registered", GetUserRegisteredBookings)
		protected.POST("/appointments/book/registered", BookRegisteredUserAppointment)
	}
}
