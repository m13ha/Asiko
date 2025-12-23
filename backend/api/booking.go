package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	apierrors "github.com/m13ha/asiko/errors/apierrors"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/requests"
)

// @Summary Get user's registered bookings
// @Description Retrieves a paginated list of all bookings made by the currently authenticated user.
// @Tags Bookings
// @Produce  application/json
// @Security BearerAuth
// @Param status query []string false "Filter by booking status (active, pending, cancelled, etc.)"
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Page size (default: 10)"
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 401 {object} responses.APIErrorResponse "Unauthorized"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /appointments/registered [get]
// @ID getUserRegisteredBookings
func (h *Handler) GetUserRegisteredBookings(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		apierrors.UnauthorizedError(c, "Unauthorized")
		return
	}

	ctx := c.Request.Context()
	rawStatuses := c.QueryArray("status")
	if len(rawStatuses) == 0 {
		if s := c.Query("status"); s != "" {
			rawStatuses = []string{s}
		}
	}

	var statuses []string
	seen := make(map[string]bool)

	for _, entry := range rawStatuses {
		tokens := []string{entry}
		if strings.Contains(entry, ",") {
			tokens = strings.Split(entry, ",")
		}

		for _, token := range tokens {
			s := strings.ToLower(strings.TrimSpace(token))
			if s != "" && !seen[s] {
				seen[s] = true
				statuses = append(statuses, s)
			}
		}
	}

	bookings, err := h.bookingService.GetUserBookings(ctx, c.Request, userID.String(), statuses)
	if err != nil {
		apierrors.InternalServerError(c, "Internal server error")
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
// @Failure 400 {object} responses.APIErrorResponse "Invalid request payload or validation error"
// @Failure 409 {object} responses.APIErrorResponse "Slot unavailable or capacity exceeded"
// @Router /appointments/book [post]
// @ID bookGuestAppointment
func (h *Handler) BookGuestAppointment(c *gin.Context) {
	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierrors.BadRequestError(c, "Invalid request payload")
		return
	}

	booking, err := h.bookingService.BookAppointment(req, "")
	if err != nil {
		apierrors.HandleAppError(c, err)
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
// @Failure 400 {object} responses.APIErrorResponse "Invalid request payload or validation error"
// @Failure 401 {object} responses.APIErrorResponse "Unauthorized"
// @Failure 409 {object} responses.APIErrorResponse "Slot unavailable or capacity exceeded"
// @Router /appointments/book/registered [post]
// @ID bookRegisteredUserAppointment
func (h *Handler) BookRegisteredUserAppointment(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		apierrors.UnauthorizedError(c, "Unauthorized")
		return
	}

	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierrors.BadRequestError(c, "Invalid request payload")
		return
	}

	booking, err := h.bookingService.BookAppointment(req, userID.String())
	if err != nil {
		apierrors.HandleAppError(c, err)
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// @Summary Get available slots for an appointment
// @Description Retrieves a paginated list of all available booking slots for a given appointment.
// @Tags Bookings
// @Produce  application/json
// @Param   app_code  path   string  true  "Appointment identifier (app_code)"
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Page size (default: 500)"
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 400 {object} responses.APIErrorResponse "Missing appointment code parameter"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /appointments/slots/{app_code} [get]
// @ID getAvailableSlots
func (h *Handler) GetAvailableSlots(c *gin.Context) {
	appcode := c.Param("app_code")
	if appcode == "" {
		apierrors.BadRequestError(c, "Missing appointment code parameter")
		return
	}

	slots, err := h.bookingService.GetAvailableSlots(c.Request, appcode)
	if err != nil {
		apierrors.InternalServerError(c, "Internal server error")
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
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Page size (default: 200)"
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 400 {object} responses.APIErrorResponse "Missing or invalid parameters"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /appointments/slots/{app_code}/by-day [get]
// @ID getAvailableSlotsByDay
func (h *Handler) GetAvailableSlotsByDay(c *gin.Context) {
	appcode := c.Param("app_code")
	if appcode == "" {
		apierrors.BadRequestError(c, "Missing appointment code parameter")
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		apierrors.BadRequestError(c, "Missing date parameter")
		return
	}

	slots, err := h.bookingService.GetAvailableSlotsByDay(c.Request, appcode, dateStr)
	if err != nil {
		apierrors.InternalServerError(c, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, slots)
}

// @Summary Get available dates for an appointment
// @Description Retrieves a list of dates that have at least one available slot.
// @Tags Bookings
// @Produce  application/json
// @Param   app_code    path   string  true  "Appointment identifier (app_code)"
// @Success 200 {array} string
// @Failure 400 {object} responses.APIErrorResponse "Missing appointment code parameter"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /appointments/dates/{app_code} [get]
// @ID getAvailableDates
func (h *Handler) GetAvailableDates(c *gin.Context) {
	appcode := c.Param("app_code")
	if appcode == "" {
		apierrors.BadRequestError(c, "Missing appointment code parameter")
		return
	}

	dates, err := h.bookingService.GetAvailableDates(c.Request.Context(), appcode)
	if err != nil {
		apierrors.InternalServerError(c, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, dates)
}

// @Summary Get all bookings for an appointment
// @Description Retrieves a paginated list of all users/bookings for a specific appointment.
// @Tags Appointments
// @Produce  application/json
// @Param   app_code  path   string  true  "Appointment identifier (app_code)"
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Page size (default: 10)"
// @Security BearerAuth
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Booking}
// @Failure 400 {object} responses.APIErrorResponse "Missing appointment code parameter"
// @Failure 401 {object} responses.APIErrorResponse "Unauthorized"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /appointments/users/{app_code} [get]
// @ID getUsersRegisteredForAppointment
func (h *Handler) GetUsersRegisteredForAppointment(c *gin.Context) {
	appCode := c.Param("app_code")
	if appCode == "" {
		apierrors.BadRequestError(c, "Missing appointment code parameter")
		return
	}

	// userID, ok := middleware.GetUUIDFromContext(c)
	// if !ok {
	// 	apierrors.UnauthorizedError(c, "Unauthorized")
	// 	return
	// }

	ctx := c.Request.Context()
	bookings, err := h.bookingService.GetAllBookingsForAppointment(ctx, c.Request, appCode)
	if err != nil {
		apierrors.InternalServerError(c, "Internal server error")
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
// @Failure 400 {object} responses.APIErrorResponse "Missing booking_code parameter"
// @Failure 404 {object} responses.APIErrorResponse "Booking not found"
// @Router /bookings/{booking_code} [get]
// @ID getBookingByCode
func (h *Handler) GetBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		apierrors.BadRequestError(c, "Missing booking_code parameter")
		return
	}

	booking, err := h.bookingService.GetBookingByCode(code)
	if err != nil {
		apierrors.NotFoundError(c, "Booking not found")
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
// @Failure 400 {object} responses.APIErrorResponse "Invalid request, validation error, or slot not available"
// @Failure 404 {object} responses.APIErrorResponse "Booking not found"
// @Failure 409 {object} responses.APIErrorResponse "Requested slot not available or capacity exceeded"
// @Router /bookings/{booking_code} [put]
// @ID updateBookingByCode
func (h *Handler) UpdateBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		apierrors.BadRequestError(c, "Missing booking_code parameter")
		return
	}

	var req requests.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierrors.BadRequestError(c, "Invalid request payload")
		return
	}

	// Validate the request using the BookingRequest's Validate method
	if err := req.Validate(); err != nil {
		apierrors.ValidationError(c, "Validation failed")
		return
	}

	booking, err := h.bookingService.UpdateBookingByCode(code, req)
	if err != nil {
		apierrors.HandleAppError(c, err)
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
// @Failure 400 {object} responses.APIErrorResponse "Error while cancelling booking"
// @Failure 404 {object} responses.APIErrorResponse "Booking not found"
// @Router /bookings/{booking_code} [delete]
// @ID cancelBookingByCode
func (h *Handler) CancelBookingByCodeHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		apierrors.BadRequestError(c, "Missing booking_code parameter")
		return
	}

	booking, err := h.bookingService.CancelBookingByCode(code)
	if err != nil {
		apierrors.HandleAppError(c, err)
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
// @Failure 400 {object} responses.APIErrorResponse "Error while rejecting booking"
// @Failure 401 {object} responses.APIErrorResponse "Unauthorized"
// @Failure 404 {object} responses.APIErrorResponse "Booking not found"
// @Router /bookings/{booking_code}/reject [post]
// @ID rejectBookingByCode
func (h *Handler) RejectBookingHandler(c *gin.Context) {
	code := c.Param("booking_code")
	if code == "" {
		apierrors.BadRequestError(c, "Missing booking_code parameter")
		return
	}

	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		apierrors.UnauthorizedError(c, "Unauthorized")
		return
	}

	booking, err := h.bookingService.RejectBooking(code, userID)
	if err != nil {
		apierrors.HandleAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, booking)
}
