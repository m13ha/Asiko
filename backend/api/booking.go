package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m13ha/asiko/errors"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/requests"
)

// @Summary Get user's registered bookings
// @Description Retrieves a paginated list of all bookings made by the currently authenticated user.
// @Tags Bookings
// @Produce  application/json
// @Security BearerAuth
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 401 {object} errors.APIErrorResponse "Unauthorized"
// @Failure 500 {object} errors.APIErrorResponse "Internal server error"
// @Router /appointments/registered [get]
// @ID getUserRegisteredBookings
func (h *Handler) GetUserRegisteredBookings(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		c.Error(errors.New(errors.CodeUnauthorized).WithKind(errors.KindUnauthorized).WithHTTP(401).WithMessage("Unauthorized"))
		return
	}

	ctx := c.Request.Context()
	bookings, err := h.bookingService.GetUserBookings(ctx, userID.String())
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// @Summary Book an appointment (Guest)
// @Description Creates a booking for an appointment as a guest user. Name and email/phone are required.
// @Tags Bookings
// @Accept  application/json
// @Produce  application/json
// @Param   booking  body   requests.BookingRequest  true  "Booking Details"
// @Success 201 {object} entities.Booking
// @Failure 400 {object} errors.APIErrorResponse "Invalid request payload or validation error"
// @Failure 409 {object} errors.APIErrorResponse "Slot unavailable or capacity exceeded"
// @Router /appointments/book [post]
// @ID bookGuestAppointment
func (h *Handler) BookGuestAppointment(c *gin.Context) {
	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request payload"))
		return
	}

	booking, err := h.bookingService.BookAppointment(req, "")
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// @Summary Book an appointment (Registered User)
// @Description Creates a booking for an appointment as a registered user.
// @Tags Bookings
// @Accept  application/json
// @Produce  application/json
// @Param   booking  body   requests.BookingRequest  true  "Booking Details"
// @Security BearerAuth
// @Success 201 {object} entities.Booking
// @Failure 400 {object} errors.APIErrorResponse "Invalid request payload or validation error"
// @Failure 401 {object} errors.APIErrorResponse "Unauthorized"
// @Failure 409 {object} errors.APIErrorResponse "Slot unavailable or capacity exceeded"
// @Router /appointments/book/registered [post]
// @ID bookRegisteredUserAppointment
func (h *Handler) BookRegisteredUserAppointment(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		c.Error(errors.New(errors.CodeUnauthorized).WithKind(errors.KindUnauthorized).WithHTTP(401).WithMessage("Unauthorized"))
		return
	}

	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request payload"))
		return
	}

	booking, err := h.bookingService.BookAppointment(req, userID.String())
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// @Summary Get available slots for an appointment
// @Description Retrieves a paginated list of all available booking slots for a given appointment.
// @Tags Bookings
// @Produce  application/json
// @Param   app_code  path   string  true  "Appointment identifier (app_code)"
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 400 {object} errors.APIErrorResponse "Missing appointment code parameter"
// @Failure 500 {object} errors.APIErrorResponse "Internal server error"
// @Router /appointments/slots/{app_code} [get]
// @ID getAvailableSlots
func (h *Handler) GetAvailableSlots(c *gin.Context) {
	appcode := c.Param("app_code")
	if appcode == "" {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Missing appointment code parameter").WithFields(errors.Field("app_code", "is required", "required")))
		return
	}

	slots, err := h.bookingService.GetAvailableSlots(c.Request, appcode)
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusOK, slots)
}

// @Summary Get available slots for a specific day
// @Description Retrieves a paginated list of available slots for an appointment on a specific day.
// @Tags Bookings
// @Produce  application/json
// @Param   app_code    path   string  true  "Appointment identifier (app_code)"
// @Param   date  query  string  true  "Date in YYYY-MM-DD format"
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 400 {object} errors.APIErrorResponse "Missing or invalid parameters"
// @Failure 500 {object} errors.APIErrorResponse "Internal server error"
// @Router /appointments/slots/{app_code}/by-day [get]
// @ID getAvailableSlotsByDay
func (h *Handler) GetAvailableSlotsByDay(c *gin.Context) {
	appcode := c.Param("app_code")
	if appcode == "" {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Missing appointment code parameter").WithFields(errors.Field("app_code", "is required", "required")))
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Missing date parameter").WithFields(errors.Field("date", "is required", "required")))
		return
	}

	slots, err := h.bookingService.GetAvailableSlotsByDay(c.Request, appcode, dateStr)
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusOK, slots)
}

// @Summary Get all bookings for an appointment
// @Description Retrieves a paginated list of all users/bookings for a specific appointment.
// @Tags Appointments
// @Produce  application/json
// @Param   app_code  path   string  true  "Appointment identifier (app_code)"
// @Security BearerAuth
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 400 {object} errors.APIErrorResponse "Missing appointment code parameter"
// @Failure 401 {object} errors.APIErrorResponse "Unauthorized"
// @Failure 500 {object} errors.APIErrorResponse "Internal server error"
// @Router /appointments/users/{app_code} [get]
// @ID getUsersRegisteredForAppointment
func (h *Handler) GetUsersRegisteredForAppointment(c *gin.Context) {
	appCode := c.Param("app_code")
	if appCode == "" {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Missing appointment code parameter").WithFields(errors.Field("app_code", "is required", "required")))
		return
	}

	ctx := c.Request.Context()
	bookings, err := h.bookingService.GetAllBookingsForAppointment(ctx, appCode)
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// @Summary Get booking by code
// @Description Retrieves booking details by its unique booking_code.
// @Tags Bookings
// @Produce  application/json
// @Param   booking_code  path   string  true  "Unique Booking Code"
// @Success 200 {object} entities.Booking
// @Failure 400 {object} errors.APIErrorResponse "Missing booking_code parameter"
// @Failure 404 {object} errors.APIErrorResponse "Booking not found"
// @Router /bookings/{booking_code} [get]
// @ID getBookingByCode
func (h *Handler) GetBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Missing booking_code parameter").WithFields(errors.Field("booking_code", "is required", "required")))
		return
	}

	booking, err := h.bookingService.GetBookingByCode(code)
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusOK, booking)
}

// @Summary Update/Reschedule a booking
// @Description Updates a booking by its unique booking_code. Can be used to reschedule.
// @Tags Bookings
// @Accept  application/json
// @Produce  application/json
// @Param   booking_code  path   string  true  "Unique Booking Code"
// @Param   booking      body   requests.BookingRequest  true  "New Booking Details"
// @Success 200 {object} entities.Booking
// @Failure 400 {object} errors.APIErrorResponse "Invalid request, validation error, or slot not available"
// @Failure 404 {object} errors.APIErrorResponse "Booking not found"
// @Failure 409 {object} errors.APIErrorResponse "Requested slot not available or capacity exceeded"
// @Router /bookings/{booking_code} [put]
// @ID updateBookingByCode
func (h *Handler) UpdateBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Missing booking_code parameter").WithFields(errors.Field("booking_code", "is required", "required")))
		return
	}

	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request payload"))
		return
	}

	// Validate the request using the BookingRequest's Validate method
	if err := req.Validate(); err != nil {
		c.Error(errors.FromValidation(err, "Validation failed"))
		return
	}

	booking, err := h.bookingService.UpdateBookingByCode(code, req)
	if err != nil {
		// Service now returns typed AppError (404/409/400); let middleware map it.
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusOK, booking)
}

// @Summary Cancel a booking
// @Description Cancels a booking by its unique booking_code. This is a soft delete.
// @Tags Bookings
// @Produce  application/json
// @Param   booking_code  path   string  true  "Unique Booking Code"
// @Success 200 {object} entities.Booking
// @Failure 400 {object} errors.APIErrorResponse "Error while cancelling booking"
// @Failure 404 {object} errors.APIErrorResponse "Booking not found"
// @Router /bookings/{booking_code} [delete]
// @ID cancelBookingByCode
func (h *Handler) CancelBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Missing booking_code parameter").WithFields(errors.Field("booking_code", "is required", "required")))
		return
	}

	booking, err := h.bookingService.CancelBookingByCode(code)
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusOK, booking)
}

// @Summary Reject a booking
// @Description Rejects a booking by its unique booking_code. This is a soft delete.
// @Tags Bookings
// @Produce  application/json
// @Param   booking_code  path   string  true  "Unique Booking Code"
// @Security BearerAuth
// @Success 200 {object} entities.Booking
// @Failure 400 {object} errors.APIErrorResponse "Error while rejecting booking"
// @Failure 401 {object} errors.APIErrorResponse "Unauthorized"
// @Failure 404 {object} errors.APIErrorResponse "Booking not found"
// @Router /bookings/{booking_code}/reject [post]
// @ID rejectBookingByCode
func (h *Handler) RejectBookingHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Missing booking_code parameter").WithFields(errors.Field("booking_code", "is required", "required")))
		return
	}

	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		c.Error(errors.New(errors.CodeUnauthorized).WithKind(errors.KindUnauthorized).WithHTTP(401).WithMessage("Unauthorized"))
		return
	}

	booking, err := h.bookingService.RejectBooking(code, userID)
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusOK, booking)
}
