package apierrors

import (
	"github.com/gin-gonic/gin"
	apperrors "github.com/m13ha/asiko/errors"
	"github.com/m13ha/asiko/models/responses"
)

// handleError sends an error response with the specified HTTP status, error code and message
func handleError(c *gin.Context, httpStatus int, errorCode string, message string) {
	c.JSON(httpStatus, responses.APIErrorResponse{
		Code:    errorCode,
		Message: message,
		HTTP:    httpStatus,
	})
}

func UnauthorizedError(c *gin.Context, message string) {
	handleError(c, 401, apperrors.CodeUnauthorized, message)
}

func ForbiddenError(c *gin.Context, message string) {
	handleError(c, 403, apperrors.CodeForbidden, message)
}

func NotFoundError(c *gin.Context, message string) {
	handleError(c, 404, apperrors.CodeResourceNotFound, message)
}

func InternalServerError(c *gin.Context, message string) {
	handleError(c, 500, apperrors.CodeInternalError, message)
}

func BadRequestError(c *gin.Context, message string) {
	handleError(c, 400, apperrors.CodeBadRequest, message)
}

func ConflictError(c *gin.Context, message string) {
	handleError(c, 409, apperrors.CodeConflict, message)
}

func ValidationError(c *gin.Context, message string) {
	handleError(c, 422, apperrors.CodeValidationFailed, message)
}

// HandleAppError translates application-specific errors into appropriate HTTP API responses.
func HandleAppError(c *gin.Context, err error) {
	appErr := apperrors.FromAppError(err)

	switch appErr.Code {
	case apperrors.CodeUnauthorized:
		UnauthorizedError(c, appErr.Message)
	case apperrors.CodeForbidden:
		ForbiddenError(c, appErr.Message)
	case apperrors.CodeResourceNotFound,
		apperrors.CodeAppointmentNotFound,
		apperrors.CodeBookingNotFound,
		apperrors.CodeRepoNotFoundError:
		NotFoundError(c, appErr.Message)
	case apperrors.CodeConflict,
		apperrors.CodeEmailAlreadyRegistered,
		apperrors.CodeBookingSlotUnavailable,
		apperrors.CodeBookingCapacityExceeded,
		apperrors.CodeDBUniqueViolation,
		apperrors.CodeRepoConflictError:
		ConflictError(c, appErr.Message)
	case apperrors.CodeValidationFailed,
		apperrors.CodeInvalidVerificationCode,
		apperrors.CodeVerificationExpired,
		apperrors.CodeRepoValidationError:
		ValidationError(c, appErr.Message)
	case apperrors.CodeLoginInvalidCredentials,
		apperrors.CodeUserPendingVerification,
		apperrors.CodeBanListBlocked:
		BadRequestError(c, appErr.Message)
	default:
		InternalServerError(c, "Internal server error")
	}
}

