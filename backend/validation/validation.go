package validation

import (
	"github.com/go-playground/validator/v10"
)

var V *validator.Validate

func InitValidator() {
	V = validator.New(validator.WithRequiredStructEnabled())
}