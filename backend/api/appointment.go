package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/m13ha/appointment_master/errors"
    "github.com/m13ha/appointment_master/middleware"
    "github.com/m13ha/appointment_master/models/requests"
    "github.com/m13ha/appointment_master/utils"
)

// parseAndValidateRequest parses and validates the appointment request from the HTTP request
func parseAndValidateRequest(c *gin.Context) (requests.AppointmentRequest, error) {
	var req requests.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, errors.NewUserError("invalid request payload: " + err.Error())
	}

	if err := utils.Validate(req); err != nil {
		return req, errors.NewUserError("validation failed: " + err.Error())
	}

	return req, nil
}

// @Summary Create a new appointment
// @Description Create a new appointment. Type can be single, group, or party.
// @Tags Appointments
// @Accept  json
// @Produce  json
// @Param   appointment  body   requests.AppointmentRequest  true  "Appointment Details"
// @Security BearerAuth
// @Success 201 {object} entities.Appointment
// @Failure 400 {object} errors.ApiErrorResponse "Invalid request payload or validation error"
// @Failure 401 {object} errors.ApiErrorResponse "Authentication required"
// @Failure 500 {object} errors.ApiErrorResponse "Failed to create appointment"
// @Router /appointments [post]
// @ID createAppointment
func (h *Handler) CreateAppointment(c *gin.Context) {
    userID, ok := middleware.GetUUIDFromContext(c)
    if !ok {
        errors.Unauthorized(c.Writer, "authentication required")
        return
    }

	req, err := parseAndValidateRequest(c)
	if err != nil {
		errors.BadRequest(c.Writer, err.Error())
		return
	}

    appointment, err := h.appointmentService.CreateAppointment(req, userID)
    if err != nil {
        switch err.(type) {
        case *errors.UserError:
            errors.BadRequest(c.Writer, err.Error())
        default:
            errors.InternalServerError(c.Writer, "failed to create appointment")
        }
        return
    }

	c.JSON(http.StatusCreated, appointment)
}

// @Summary Get appointments created by the user
// @Description Retrieves a paginated list of appointments created by the currently authenticated user.
// @Tags Appointments
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} responses.PaginatedResponse{items=[]responses.AppointmentResponse}
// @Failure 401 {object} errors.ApiErrorResponse "Authentication required"
// @Router /appointments/my [get]
// @ID getMyAppointments
func (h *Handler) GetAppointmentsCreatedByUser(c *gin.Context) {
    uid, ok := middleware.GetUUIDFromContext(c)
    if !ok {
        errors.Unauthorized(c.Writer, "")
        return
    }

    appointments := h.appointmentService.GetAllAppointmentsCreatedByUser(uid.String(), nil)

	c.JSON(http.StatusOK, appointments)
}
