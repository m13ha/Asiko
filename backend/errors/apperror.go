package errors

import (
    "errors"
    "fmt"
    "strings"
)

// Kind is a high-level error category used for HTTP/status mapping and logging.
type Kind string

const (
    KindValidation    Kind = "validation"
    KindUnauthorized  Kind = "unauthorized"
    KindForbidden     Kind = "forbidden"
    KindNotFound      Kind = "not_found"
    KindConflict      Kind = "conflict"
    KindRateLimited   Kind = "rate_limited"
    KindPrecondition  Kind = "precondition_failed"
    KindTimeout       Kind = "timeout"
    KindCanceled      Kind = "canceled"
    KindExternal      Kind = "external"
    KindInternal      Kind = "internal"
)

// FieldError describes a validation error for a specific field.
type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Rule    string `json:"rule,omitempty"`
}

// AppError is the canonical error type across the layers.
type AppError struct {
    Code    string            // stable machine code for FE handling
    Kind    Kind              // category for mapping
    HTTP    int               // preferred HTTP status (optional override)
    Message string            // safe user-facing message
    Fields  []FieldError      // validation fields, if any
    Op      string            // operation (package.func)
    Meta    map[string]any    // safe structured context
    Cause   error             // underlying error
}

func (e *AppError) Error() string {
    if e == nil {
        return "<nil>"
    }
    if e.Code != "" && e.Message != "" {
        return fmt.Sprintf("%s: %s", e.Code, e.Message)
    }
    if e.Message != "" {
        return e.Message
    }
    return e.Code
}

func (e *AppError) Unwrap() error { return e.Cause }

// New constructs a minimal AppError.
func New(code string) *AppError { return &AppError{Code: code} }

func (e *AppError) WithKind(k Kind) *AppError           { e.Kind = k; return e }
func (e *AppError) WithHTTP(code int) *AppError          { e.HTTP = code; return e }
func (e *AppError) WithMessage(msg string) *AppError     { e.Message = msg; return e }
func (e *AppError) WithFields(fs ...FieldError) *AppError { e.Fields = append(e.Fields, fs...); return e }
func (e *AppError) WithOp(op string) *AppError           { e.Op = op; return e }
func (e *AppError) WithMeta(k string, v any) *AppError   { if e.Meta == nil { e.Meta = map[string]any{} }; e.Meta[k] = v; return e }
func (e *AppError) WithCause(err error) *AppError        { e.Cause = err; return e }

// FromError converts any error into an AppError, preserving an existing AppError.
func FromError(err error) *AppError {
    if err == nil {
        return nil
    }
    var ae *AppError
    if errors.As(err, &ae) {
        return ae
    }
    // Map legacy types
    if pe, ok := err.(*PendingVerificationError); ok {
        return New(CodeUserPendingVerification).WithKind(KindPrecondition).WithHTTP(202).WithMessage(pe.Message).WithCause(err)
    }
    if ue, ok := err.(*UserError); ok {
        // Heuristic mapping based on message for better codes
        msg := ue.Message
        lower := strings.ToLower(msg)
        switch {
        case strings.Contains(lower, "pending") && strings.Contains(lower, "verification"):
            return New(CodeUserPendingVerification).WithKind(KindPrecondition).WithHTTP(202).WithMessage(msg).WithCause(err)
        case strings.Contains(lower, "invalid email or password"):
            return New(CodeLoginInvalidCredentials).WithKind(KindUnauthorized).WithHTTP(401).WithMessage("Invalid email or password").WithCause(err)
        case strings.Contains(lower, "email already registered"):
            return New(CodeEmailAlreadyRegistered).WithKind(KindConflict).WithHTTP(409).WithMessage("Email already registered").WithCause(err)
        case strings.Contains(lower, "appointment not found"):
            return New(CodeAppointmentNotFound).WithKind(KindNotFound).WithHTTP(404).WithMessage("Appointment not found").WithCause(err)
        case strings.Contains(lower, "booking not found"):
            return New(CodeBookingNotFound).WithKind(KindNotFound).WithHTTP(404).WithMessage("Booking not found").WithCause(err)
        case strings.Contains(lower, "no available slot"):
            return New(CodeBookingSlotUnavailable).WithKind(KindConflict).WithHTTP(409).WithMessage("No available slot").WithCause(err)
        case strings.Contains(lower, "not enough capacity"):
            return New(CodeBookingCapacityExceeded).WithKind(KindConflict).WithHTTP(409).WithMessage("Not enough capacity").WithCause(err)
        case strings.Contains(lower, "forbidden") || strings.Contains(lower, "unauthorized"):
            return New(CodeUnauthorized).WithKind(KindUnauthorized).WithHTTP(401).WithMessage(msg).WithCause(err)
        default:
            return New(CodeValidationFailed).WithKind(KindValidation).WithHTTP(400).WithMessage(msg).WithCause(err)
        }
    }
    // Fallback internal
    return New(CodeInternalError).WithKind(KindInternal).WithHTTP(500).WithMessage("Internal server error").WithCause(err)
}
