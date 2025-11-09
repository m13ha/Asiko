package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/m13ha/asiko/errors"
    "github.com/m13ha/asiko/middleware"
    "github.com/m13ha/asiko/models/requests"
    "github.com/m13ha/asiko/utils"
)

// parseAndValidateRequest parses and validates the appointment request from the HTTP request
func parseAndValidateRequest(c *gin.Context) (requests.AppointmentRequest, error) {
    var req requests.AppointmentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        return req, errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("invalid request payload: "+err.Error())
    }

    if err := utils.Validate(req); err != nil {
        return req, errors.FromValidation(err, "validation failed")
    }

    return req, nil
}

// @Summary Create a new appointment
// @Description Create a new appointment. Type can be single, group, or party.
// @Tags Appointments
// @Accept  application/json
// @Produce  application/json
// @Param   appointment  body   requests.AppointmentRequest  true  "Appointment Details"
// @Security BearerAuth
// @Success 201 {object} entities.Appointment
// @Failure 400 {object} errors.APIErrorResponse "Invalid request payload or validation error"
// @Failure 401 {object} errors.APIErrorResponse "Authentication required"
// @Failure 500 {object} errors.APIErrorResponse "Failed to create appointment"
// @Router /appointments [post]
// @ID createAppointment
func (h *Handler) CreateAppointment(c *gin.Context) {
    userID, ok := middleware.GetUUIDFromContext(c)
    if !ok {
        c.Error(errors.New(errors.CodeUnauthorized).WithKind(errors.KindUnauthorized).WithHTTP(401).WithMessage("authentication required"))
        return
    }

    req, err := parseAndValidateRequest(c)
    if err != nil {
        c.Error(errors.FromError(err))
        return
    }

    appointment, err := h.appointmentService.CreateAppointment(req, userID)
    if err != nil {
        c.Error(errors.FromError(err))
        return
    }

    c.JSON(http.StatusCreated, appointment)
}

// @Summary Get appointments created by the user
// @Description Retrieves a paginated list of appointments created by the currently authenticated user.
// @Tags Appointments
// @Produce  application/json
// @Security BearerAuth
// @Success 200 {object} responses.PaginatedResponse{items=[]responses.AppointmentResponse}
// @Failure 401 {object} errors.APIErrorResponse "Authentication required"
// @Router /appointments/my [get]
// @ID getMyAppointments
func (h *Handler) GetAppointmentsCreatedByUser(c *gin.Context) {
    uid, ok := middleware.GetUUIDFromContext(c)
    if !ok {
        c.Error(errors.New(errors.CodeUnauthorized).WithKind(errors.KindUnauthorized).WithHTTP(401).WithMessage("Authentication required"))
        return
    }

    appointments := h.appointmentService.GetAllAppointmentsCreatedByUser(uid.String(), c.Request)

    c.JSON(http.StatusOK, appointments)
}
