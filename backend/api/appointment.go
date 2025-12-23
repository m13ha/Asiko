package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	apierrors "github.com/m13ha/asiko/errors/apierrors"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	"github.com/m13ha/asiko/utils"
)

func parseStatusFilters(raw []string) []entities.AppointmentStatus {
	values := make([]entities.AppointmentStatus, 0, len(raw))
	seen := map[entities.AppointmentStatus]struct{}{}

	normalize := func(token string) {
		token = strings.TrimSpace(strings.ToLower(token))
		switch token {
		case string(entities.AppointmentStatusPending),
			string(entities.AppointmentStatusOngoing),
			string(entities.AppointmentStatusCompleted),
			string(entities.AppointmentStatusCanceled),
			string(entities.AppointmentStatusExpired):
			status := entities.AppointmentStatus(token)
			if _, ok := seen[status]; !ok {
				seen[status] = struct{}{}
				values = append(values, status)
			}
		}
	}

	for _, entry := range raw {
		if strings.Contains(entry, ",") {
			for _, token := range strings.Split(entry, ",") {
				normalize(token)
			}
			continue
		}
		normalize(entry)
	}

	return values
}

// @Summary Create a new appointment
// @Description Create a new appointment. Type can be single, group, or party.
// @Tags Appointments
// @Accept  application/json
// @Produce  application/json
// @Param   appointment  body   requests.AppointmentRequest  true  "Appointment Details"
// @Security BearerAuth
// @Success 201 {object} entities.Appointment
// @Failure 400 {object} responses.APIErrorResponse "Invalid request payload or validation error"
// @Failure 401 {object} responses.APIErrorResponse "Authentication required"
// @Failure 500 {object} responses.APIErrorResponse "Failed to create appointment"
// @Router /appointments [post]
// @ID createAppointment
func (h *Handler) CreateAppointment(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		apierrors.UnauthorizedError(c, "authentication required")
		return
	}

	// Parse and validate request manually to use the new error system
	var req requests.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierrors.BadRequestError(c, "invalid request payload: "+err.Error())
		return
	}

	if err := utils.Validate(req); err != nil {
		apierrors.ValidationError(c, "validation failed")
		return
	}

	appointment, err := h.appointmentService.CreateAppointment(req, userID)
	if err != nil {
		apierrors.InternalServerError(c, "Failed to create appointment")
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// @Summary Get appointments created by the user
// @Description Retrieves a paginated list of appointments created by the currently authenticated user.
// @Tags Appointments
// @Produce  application/json
// @Security BearerAuth
// @Param status query []string false "Filter by appointment status (pending, ongoing, completed, canceled, expired)"
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Page size (default: 10)"
// @Success 200 {object} responses.PaginatedResponse{items=[]responses.AppointmentResponse}
// @Failure 401 {object} responses.APIErrorResponse "Authentication required"
// @Router /appointments/my [get]
// @ID getMyAppointments
func (h *Handler) GetAppointmentsCreatedByUser(c *gin.Context) {
	uid, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		apierrors.UnauthorizedError(c, "Authentication required")
		return
	}

	statuses := parseStatusFilters(c.QueryArray("status"))
	appointments := h.appointmentService.GetAllAppointmentsCreatedByUser(uid.String(), c.Request, statuses)

	c.JSON(http.StatusOK, appointments)
}

// @Summary Get appointment by app code
// @Description Retrieves appointment details by its unique app_code, public endpoint for booking flow
// @Tags Appointments
// @Produce  application/json
// @Param   app_code  path   string  true  "Appointment identifier (app_code)"
// @Success 200 {object} entities.Appointment
// @Failure 400 {object} responses.APIErrorResponse "Missing app_code parameter"
// @Failure 404 {object} responses.APIErrorResponse "Appointment not found"
// @Router /appointments/code/{app_code} [get]
// @ID getAppointmentByAppCode
func (h *Handler) GetAppointmentByAppCode(c *gin.Context) {
	appCode := c.Param("app_code")
	if appCode == "" {
		apierrors.BadRequestError(c, "Missing app_code parameter")
		return
	}

	appointment, err := h.appointmentService.GetAppointmentByAppCode(appCode)
	if err != nil {
		apierrors.NotFoundError(c, "Appointment not found")
		return
	}

	c.JSON(http.StatusOK, appointment)
}
