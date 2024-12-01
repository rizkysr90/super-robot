package helper

import "github.com/go-playground/validator/v10"

func FormatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "max":
		return e.Field() + " exceeds maximum length"
	case "gte":
		return e.Field() + " must be greater than or equal to 0"
	case "uuid":
		return e.Field() + " must be a valid UUID"
	default:
		return e.Field() + " is invalid"
	}
}
