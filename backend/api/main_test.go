package api

import (
	"github.com/gin-gonic/gin"
	"github.com/m13ha/appointment_master/services/mocks"
)

// setupTestRouter initializes a gin router with mocked services for testing.
func setupTestRouter() (*gin.Engine, *mocks.UserService, *mocks.AppointmentService, *mocks.BookingService, *mocks.AnalyticsService, *mocks.BanListService, *mocks.EventNotificationService) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(mocks.UserService)
	mockAppointmentService := new(mocks.AppointmentService)
	mockBookingService := new(mocks.BookingService)
	mockAnalyticsService := new(mocks.AnalyticsService)
	mockBanListService := new(mocks.BanListService)
	mockEventNotificationService := new(mocks.EventNotificationService)

	h := NewHandler(mockUserService, mockAppointmentService, mockBookingService, mockAnalyticsService, mockBanListService, mockEventNotificationService)

	router := gin.Default()
	router.POST("/login", h.Login)
	router.POST("/users", h.CreateUser)
	router.POST("/auth/verify-registration", h.VerifyRegistrationHandler)
	router.POST("/appointments", h.CreateAppointment)
	router.DELETE("/bookings/:booking_code", h.CancelBookingByCodeHandler)
	router.GET("/analytics", h.GetUserAnalytics)

	return router, mockUserService, mockAppointmentService, mockBookingService, mockAnalyticsService, mockBanListService, mockEventNotificationService
}
