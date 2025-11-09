package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/m13ha/asiko/errors"
    "github.com/m13ha/asiko/middleware"
    "github.com/m13ha/asiko/models/requests"
    "github.com/m13ha/asiko/models/responses"
)

// @Summary Add email to ban list
// @Description Add an email to the user's personal ban list.
// @Tags BanList
// @Accept  application/json
// @Produce  application/json
// @Param   ban_request  body   requests.BanRequest  true  "Email to ban"
// @Security BearerAuth
// @Success 201 {object} entities.BanListEntry
// @Failure 400 {object} errors.APIErrorResponse "Invalid request"
// @Failure 401 {object} errors.APIErrorResponse "Unauthorized"
// @Failure 409 {object} errors.APIErrorResponse "Email already on ban list"
// @Router /ban-list [post]
// @ID addToBanList
func (h *Handler) AddToBanList(c *gin.Context) {
    userID, ok := middleware.GetUUIDFromContext(c)
    if !ok {
        c.Error(errors.New(errors.CodeUnauthorized).WithKind(errors.KindUnauthorized).WithHTTP(401).WithMessage("Unauthorized"))
        return
    }

    var req requests.BanRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request payload"))
        return
    }

    entry, err := h.banService.AddToBanList(userID, req.Email)
    if err != nil {
        c.Error(errors.FromError(err))
        return
    }

    c.JSON(http.StatusCreated, entry)
}

// @Summary Remove email from ban list
// @Description Remove an email from the user's personal ban list.
// @Tags BanList
// @Accept  application/json
// @Produce  application/json
// @Param   ban_request  body   requests.BanRequest  true  "Email to unban"
// @Security BearerAuth
// @Success 200 {object} responses.SimpleMessage
// @Failure 400 {object} errors.APIErrorResponse "Invalid request"
// @Failure 401 {object} errors.APIErrorResponse "Unauthorized"
// @Router /ban-list [delete]
// @ID removeFromBanList
func (h *Handler) RemoveFromBanList(c *gin.Context) {
    userID, ok := middleware.GetUUIDFromContext(c)
    if !ok {
        c.Error(errors.New(errors.CodeUnauthorized).WithKind(errors.KindUnauthorized).WithHTTP(401).WithMessage("Unauthorized"))
        return
    }

    var req requests.BanRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request payload"))
        return
    }

    if err := h.banService.RemoveFromBanList(userID, req.Email); err != nil {
        c.Error(errors.FromError(err))
        return
    }

    c.JSON(http.StatusOK, responses.SimpleMessage{Message: "Email removed from ban list"})
}

// @Summary Get user's ban list
// @Description Get a list of all emails on the user's personal ban list.
// @Tags BanList
// @Produce  application/json
// @Security BearerAuth
// @Success 200 {array} entities.BanListEntry
// @Failure 401 {object} errors.APIErrorResponse "Unauthorized"
// @Router /ban-list [get]
// @ID getBanList
func (h *Handler) GetBanList(c *gin.Context) {
    userID, ok := middleware.GetUUIDFromContext(c)
    if !ok {
        c.Error(errors.New(errors.CodeUnauthorized).WithKind(errors.KindUnauthorized).WithHTTP(401).WithMessage("Unauthorized"))
        return
    }

    banList, err := h.banService.GetBanList(userID)
    if err != nil {
        c.Error(errors.FromError(err))
        return
    }

    c.JSON(http.StatusOK, banList)
}
