package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/middleware"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/models/responses"
)

// @Summary Add email to ban list
// @Description Add an email to the user's personal ban list.
// @Tags Ban List
// @Accept  json
// @Produce  json
// @Param   ban_request  body   requests.BanRequest  true  "Email to ban"
// @Security BearerAuth
// @Success 201 {object} entities.BanListEntry
// @Failure 400 {object} errors.ApiErrorResponse "Invalid request"
// @Failure 401 {object} errors.ApiErrorResponse "Unauthorized"
// @Router /ban-list [post]
// @ID addToBanList
func (h *Handler) AddToBanList(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		errors.Unauthorized(c.Writer, "Unauthorized")
		return
	}

	var req requests.BanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request payload")
		return
	}

	entry, err := h.banService.AddToBanList(userID, req.Email)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// @Summary Remove email from ban list
// @Description Remove an email from the user's personal ban list.
// @Tags Ban List
// @Accept  json
// @Produce  json
// @Param   ban_request  body   requests.BanRequest  true  "Email to unban"
// @Security BearerAuth
// @Success 200 {object} responses.SimpleMessage
// @Failure 400 {object} errors.ApiErrorResponse "Invalid request"
// @Failure 401 {object} errors.ApiErrorResponse "Unauthorized"
// @Router /ban-list [delete]
// @ID removeFromBanList
func (h *Handler) RemoveFromBanList(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		errors.Unauthorized(c.Writer, "Unauthorized")
		return
	}

	var req requests.BanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request payload")
		return
	}

	if err := h.banService.RemoveFromBanList(userID, req.Email); err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, responses.SimpleMessage{Message: "Email removed from ban list"})
}

// @Summary Get user's ban list
// @Description Get a list of all emails on the user's personal ban list.
// @Tags Ban List
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} entities.BanListEntry
// @Failure 401 {object} errors.ApiErrorResponse "Unauthorized"
// @Router /ban-list [get]
// @ID getBanList
func (h *Handler) GetBanList(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		errors.Unauthorized(c.Writer, "Unauthorized")
		return
	}

	banList, err := h.banService.GetBanList(userID)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, banList)
}
