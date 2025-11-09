package errors

import "net/http"

// StatusFromKind maps an error Kind to a default HTTP status.
func StatusFromKind(k Kind) int {
    switch k {
    case KindValidation:
        return http.StatusBadRequest
    case KindUnauthorized:
        return http.StatusUnauthorized
    case KindForbidden:
        return http.StatusForbidden
    case KindNotFound:
        return http.StatusNotFound
    case KindConflict:
        return http.StatusConflict
    case KindRateLimited:
        return http.StatusTooManyRequests
    case KindPrecondition:
        return http.StatusPreconditionFailed
    case KindTimeout:
        return http.StatusGatewayTimeout
    case KindCanceled:
        return http.StatusRequestTimeout
    case KindExternal:
        return http.StatusBadGateway
    case KindInternal:
        fallthrough
    default:
        return http.StatusInternalServerError
    }
}

// APIErrorResponse is the standardized error response payload.
type APIErrorResponse struct {
    Status    int          `json:"status"`
    Code      string       `json:"code"`
    Message   string       `json:"message"`
    Fields    []FieldError `json:"fields,omitempty"`
    RequestID string       `json:"request_id,omitempty"`
    Meta      any          `json:"meta,omitempty"`
}

