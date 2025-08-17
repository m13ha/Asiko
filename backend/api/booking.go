package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/middleware"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/services"
	"github.com/m13ha/appointment_master/utils"
)

func GetUserRegisteredBookings(c *gin.Context) {
	userIDStr := middleware.GetUserIDFromContext(c)
	if userIDStr == "" {
		errors.Unauthorized(c.Writer, "")
		return
	}

	bookings, err := services.GetUserBookings(userIDStr, c.Request)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// BookGuestAppointment handles guest bookings
func BookGuestAppointment(c *gin.Context) {
	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request payload")
		return
	}

	if err := utils.Validate(req); err != nil {
		errors.BadRequest(c.Writer, "Validation failed: "+err.Error())
		return
	}

	bookingResponse, err := services.BookGuestAppointment(req)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, bookingResponse)
}

// BookRegisteredUserAppointment handles registered user bookings
func BookRegisteredUserAppointment(c *gin.Context) {
	userIDStr := middleware.GetUserIDFromContext(c)
	if userIDStr == "" {
		errors.Unauthorized(c.Writer, "")
		return
	}

	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request payload")
		return
	}

	if err := utils.Validate(req); err != nil {
		errors.BadRequest(c.Writer, "Validation failed: "+err.Error())
		return
	}

	bookingResponse, err := services.BookRegisteredUserAppointment(req, userIDStr)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, bookingResponse)
}

func GetAvailableSlots(c *gin.Context) {
	appcode := c.Param("id")
	if appcode == "" {
		errors.BadRequest(c.Writer, "Missing appointment code parameter")
		return
	}

	slots, err := services.GetAvailableSlots(appcode, c.Request)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, slots)
}

// GetAvailableSlotsByDay returns available slots for an appointment on a specific day
func GetAvailableSlotsByDay(c *gin.Context) {
	appcode := c.Param("id")
	if appcode == "" {
		errors.BadRequest(c.Writer, "Missing appointment code parameter")
		return
	}
	dateStr := c.Query("date")
	if dateStr == "" {
		errors.BadRequest(c.Writer, "Missing date parameter")
		return
	}
	slots, err := services.GetAvailableSlotsByDay(appcode, dateStr, c.Request)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, slots)
}

func GetUsersRegisteredForAppointment(c *gin.Context) {
	app_code := c.Param("id")
	if app_code == "" {
		errors.BadRequest(c.Writer, "Missing appointment code parameter")
		return
	}
	bookings, err := services.GetAllBookingsForAppointment(app_code, c.Request)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// GetBookingByCodeHandler returns booking details by booking_code
func GetBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		errors.BadRequest(c.Writer, "Missing booking_code parameter")
		return
	}
	booking, err := services.GetBookingByCode(code)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, booking)
}

// UpdateBookingByCodeHandler reschedules a booking if slot is available
func UpdateBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		errors.BadRequest(c.Writer, "Missing booking_code parameter")
		return
	}
	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request payload")
		return
	}
	if err := utils.Validate(req); err != nil {
		errors.BadRequest(c.Writer, "Validation failed: "+err.Error())
		return
	}
	bookingResponse, err := services.UpdateBookingByCode(code, req)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, bookingResponse)
}

// CancelBookingByCodeHandler cancels a booking by booking_code
func CancelBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		errors.BadRequest(c.Writer, "Missing booking_code parameter")
		return
	}
	bookingResponse, err := services.CancelBookingByCode(code)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, bookingResponse)
}
