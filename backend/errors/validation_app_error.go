package errors

import (
    "strings"
    "github.com/go-playground/validator/v10"
)

// Field constructs a FieldError helper.
func Field(name, message, rule string) FieldError {
    return FieldError{Field: name, Message: message, Rule: rule}
}

// FromValidation converts validator.ValidationErrors to an AppError with Kind=validation.
// If err is not a validator error, returns a generic VALIDATION_FAILED AppError with the
// provided fallback message.
func FromValidation(err error, fallbackMsg string) *AppError {
    if err == nil {
        return nil
    }
    var fields []FieldError
    if verrs, ok := err.(validator.ValidationErrors); ok {
        for _, e := range verrs {
            // Build a concise message; keep rule for clients to specialize messages.
            fields = append(fields, FieldError{
                Field:   e.Field(),
                Message: e.Error(),
                Rule:    e.Tag(),
            })
        }
        return New(CodeValidationFailed).
            WithKind(KindValidation).
            WithHTTP(400).
            WithMessage("Validation failed").
            WithFields(fields...)
    }
    // Fallback when bind/unmarshal failed but not a validator set
    msg := strings.TrimSpace(fallbackMsg)
    if msg == "" {
        msg = "Invalid request"
    }
    return New(CodeValidationFailed).WithKind(KindValidation).WithHTTP(400).WithMessage(msg)
}

