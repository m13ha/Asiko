package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// Validate validates a struct using the validator package
func Validate(i interface{}) error {
	return validate.Struct(i)
}
