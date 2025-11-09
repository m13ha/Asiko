package errors

import (
	"strings"
)

// TranslateRepoError maps repository/DB-layer errors into AppError.
// Implementors can call this before returning errors to services.
func TranslateRepoError(op string, err error) *AppError {
	if err == nil {
		return nil
	}
	// gorm.ErrRecordNotFound
	if strings.Contains(err.Error(), "record not found") {
		return New(CodeResourceNotFound).WithKind(KindNotFound).WithMessage("Resource not found").WithOp(op).WithCause(err)
	}
	// Unique violation (portable best-effort; for specific drivers match code)
	if strings.Contains(strings.ToLower(err.Error()), "unique") {
		return New(CodeDBUniqueViolation).WithKind(KindConflict).WithMessage("Conflict: duplicate value").WithOp(op).WithCause(err)
	}
	// Fallback
	return FromError(err).WithOp(op)
}
