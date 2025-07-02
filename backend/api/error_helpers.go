package api

import (
	"net/http"

	myerrors "github.com/m13ha/appointment_master/errors"
)

// HandleServiceError checks for UserError and writes a user-friendly message, otherwise writes a generic error
func HandleServiceError(w http.ResponseWriter, err error, notFoundStatus int, writeError func(http.ResponseWriter, int, string)) {
	if err == nil {
		return
	}
	if ue, ok := err.(*myerrors.UserError); ok {
		writeError(w, notFoundStatus, ue.Message)
	} else {
		writeError(w, http.StatusInternalServerError, "An unexpected error occurred. Please try again later.")
	}
}
