package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/middleware"
	"github.com/m13ha/appointment_master/models/requests"
)

// @Summary Get user's registered bookings
// @Description Retrieves a paginated list of all bookings made by the currently authenticated user.
// @Tags Bookings
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 401 {object} errors.ApiErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ApiErrorResponse "Internal server error"
// @Router /appointments/registered [get]
// @ID getUserRegisteredBookings
func (h *Handler) GetUserRegisteredBookings(c *gin.Context) {
	userIDStr := middleware.GetUserIDFromContext(c)
	if userIDStr == "" {
		errors.Unauthorized(c.Writer, "Unauthorized")
		return
	}

	ctx := c.Request.Context()
	bookings, err := h.bookingService.GetUserBookings(ctx, userIDStr)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// @Summary Book an appointment (Guest)
// @Description Creates a booking for an appointment as a guest user. Name and email/phone are required.
// @Tags Bookings
// @Accept  json
// @Produce  json
// @Param   booking  body   requests.BookingRequest  true  "Booking Details"
// @Success 201 {object} entities.Booking
// @Failure 400 {object} errors.ApiErrorResponse "Invalid request payload, validation error, or capacity exceeded"
// @Router /appointments/book [post]
// @ID bookGuestAppointment
func (h *Handler) BookGuestAppointment(c *gin.Context) {
	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request payload")
		return
	}

	booking, err := h.bookingService.BookAppointment(req, "")
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// @Summary Book an appointment (Registered User)
// @Description Creates a booking for an appointment as a registered user.
// @Tags Bookings
// @Accept  json
// @Produce  json
// @Param   booking  body   requests.BookingRequest  true  "Booking Details"
// @Security BearerAuth
// @Success 201 {object} entities.Booking
// @Failure 400 {object} errors.ApiErrorResponse "Invalid request payload, validation error, or capacity exceeded"
// @Failure 401 {object} errors.ApiErrorResponse "Unauthorized"
// @Router /appointments/book/registered [post]
// @ID bookRegisteredUserAppointment
func (h *Handler) BookRegisteredUserAppointment(c *gin.Context) {
	userIDStr := middleware.GetUserIDFromContext(c)
	if userIDStr == "" {
		errors.Unauthorized(c.Writer, "Unauthorized")
		return
	}

	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request payload")
		return
	}

	booking, err := h.bookingService.BookAppointment(req, userIDStr)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// @Summary Get available slots for an appointment
// @Description Retrieves a paginated list of all available booking slots for a given appointment.
// @Tags Bookings
// @Produce  json
// @Param   id  path   string  true  "Appointment Code (app_code)"
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 400 {object} errors.ApiErrorResponse "Missing appointment code parameter"
// @Failure 500 {object} errors.ApiErrorResponse "Internal server error"
// @Router /appointments/slots/{id} [get]
// @ID getAvailableSlots
func (h *Handler) GetAvailableSlots(c *gin.Context) {
	appcode := c.Param("id")
	if appcode == "" {
		errors.BadRequest(c.Writer, "Missing appointment code parameter")
		return
	}

	ctx := c.Request.Context()
	slots, err := h.bookingService.GetAvailableSlots(ctx, appcode)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, slots)
}

// @Summary Get available slots for a specific day
// @Description Retrieves a paginated list of available slots for an appointment on a specific day.
// @Tags Bookings
// @Produce  json
// @Param   id    path   string  true  "Appointment Code (app_code)"
// @Param   date  query  string  true  "Date in YYYY-MM-DD format"
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 400 {object} errors.ApiErrorResponse "Missing or invalid parameters"
// @Failure 500 {object} errors.ApiErrorResponse "Internal server error"
// @Router /appointments/slots/{id}/by-day [get] // Note: This route is not in handlers.go, assuming it should be added
// @ID getAvailableSlotsByDay
func (h *Handler) GetAvailableSlotsByDay(c *gin.Context) {
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

	ctx := c.Request.Context()
	slots, err := h.bookingService.GetAvailableSlotsByDay(ctx, appcode, dateStr)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, slots)
}

// @Summary Get all bookings for an appointment
// @Description Retrieves a paginated list of all users/bookings for a specific appointment.
// @Tags Appointments
// @Produce  json
// @Param   id  path   string  true  "Appointment Code (app_code)"
// @Security BearerAuth
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 400 {object} errors.ApiErrorResponse "Missing appointment code parameter"
// @Failure 401 {object} errors.ApiErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ApiErrorResponse "Internal server error"
// @Router /appointments/users/{id} [get]
// @ID getUsersRegisteredForAppointment
func (h *Handler) GetUsersRegisteredForAppointment(c *gin.Context) {
	appCode := c.Param("id")
	if appCode == "" {
		errors.BadRequest(c.Writer, "Missing appointment code parameter")
		return
	}

	ctx := c.Request.Context()
	bookings, err := h.bookingService.GetAllBookingsForAppointment(ctx, appCode)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// @Summary Get booking by code
// @Description Retrieves booking details by its unique booking_code.
// @Tags Bookings
// @Produce  json
// @Param   booking_code  path   string  true  "Unique Booking Code"
// @Success 200 {object} entities.Booking
// @Failure 400 {object} errors.ApiErrorResponse "Missing booking_code parameter"
// @Failure 404 {object} errors.ApiErrorResponse "Booking not found"
// @Router /bookings/{booking_code} [get]
// @ID getBookingByCode
func (h *Handler) GetBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		errors.BadRequest(c.Writer, "Missing booking_code parameter")
		return
	}

	booking, err := h.bookingService.GetBookingByCode(code)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, booking)
}

// @Summary Update/Reschedule a booking
// @Description Updates a booking by its unique booking_code. Can be used to reschedule.
// @Tags Bookings
// @Accept  json
// @Produce  json
// @Param   booking_code  path   string  true  "Unique Booking Code"
// @Param   booking      body   requests.BookingRequest  true  "New Booking Details"
// @Success 200 {object} entities.Booking
// @Failure 400 {object} errors.ApiErrorResponse "Invalid request, validation error, or slot not available"
// @Failure 404 {object} errors.ApiErrorResponse "Booking not found"
// @Router /bookings/{booking_code} [put]
// @ID updateBookingByCode
func (h *Handler) UpdateBookingByCodeHandler(c *gin.Context) {
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

	// Validate the request using the BookingRequest's Validate method
	if err := req.Validate(); err != nil {
		errors.BadRequest(c.Writer, "Validation failed: "+err.Error())
		return
	}

	booking, err := h.bookingService.UpdateBookingByCode(code, req)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, booking)
}

// @Summary Cancel a booking
// @Description Cancels a booking by its unique booking_code. This is a soft delete.
// @Tags Bookings
// @Produce  json
// @Param   booking_code  path   string  true  "Unique Booking Code"
// @Success 200 {object} entities.Booking
// @Failure 400 {object} errors.ApiErrorResponse "Error while cancelling booking"
// @Failure 404 {object} errors.ApiErrorResponse "Booking not found"
// @Router /bookings/{booking_code} [delete]
// @ID cancelBookingByCode
func (h *Handler) CancelBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		errors.BadRequest(c.Writer, "Missing booking_code parameter")
		return
	}

	booking, err := h.bookingService.CancelBookingByCode(code)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, booking)
}

// @Summary Reject a booking
// @Description Rejects a booking by its unique booking_code. This is a soft delete.
// @Tags Bookings
// @Produce  json
// @Param   booking_code  path   string  true  "Unique Booking Code"
// @Security BearerAuth
// @Success 200 {object} entities.Booking
// @Failure 400 {object} errors.ApiErrorResponse "Error while rejecting booking"
// @Failure 401 {object} errors.ApiErrorResponse "Unauthorized"
// @Failure 404 {object} errors.ApiErrorResponse "Booking not found"
// @Router /bookings/{booking_code}/reject [post]
// @ID rejectBookingByCode
func (h *Handler) RejectBookingHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		errors.BadRequest(c.Writer, "Missing booking_code parameter")
		return
	}

	userIDStr := middleware.GetUserIDFromContext(c)
	if userIDStr == "" {
		errors.Unauthorized(c.Writer, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		errors.BadRequest(c.Writer, "Invalid user ID")
		return
	}

	booking, err := h.bookingService.RejectBooking(code, userID)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, booking)
}
