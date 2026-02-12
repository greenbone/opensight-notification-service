package validation

import (
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New(validator.WithRequiredStructEnabled()) // singleton validator instance to be used throughout the application
