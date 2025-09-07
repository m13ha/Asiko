package api

import (
	"github.com/gin-gonic/gin"
	"github.com/m13ha/appointment_master/services/mocks"
)

// setupTestRouter initializes a gin router with mocked services for testing.
func setupTestRouter() (*gin.Engine, *mocks.UserService, *mocks.AppointmentService, *mocks.BookingService) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(mocks.UserService)
	mockAppointmentService := new(mocks.AppointmentService)
	mockBookingService := new(mocks.BookingService)

	// We pass the real NewHandler function but with our mocked services.
	h := NewHandler(mockUserService, mockAppointmentService, mockBookingService)

	router := gin.Default()
	// We are only testing the handlers, so we register routes for a specific handler.
	// This could be expanded to register all routes if needed for broader tests.
	router.POST("/login", h.Login)
	router.POST("/users", h.CreateUser)
	router.POST("/appointments", h.CreateAppointment)
	router.DELETE("/bookings/:booking_code", h.CancelBookingByCodeHandler)

	return router, mockUserService, mockAppointmentService, mockBookingService
}
