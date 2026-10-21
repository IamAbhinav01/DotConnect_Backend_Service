package validator

import "github.com/go-playground/validator/v10"

var Validate validator.Validate

func NewValidator() validator.Validate {
	return *validator.New(
		validator.WithRequiredStructEnabled(),
	)
}

func init() {
	Validate = NewValidator()
}