package utility

import (
	"rizkysr90-pos/internal/helper"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/go-playground/validator/v10"
)

// ValidationUtil holds validator instance and methods for validation
type ValidationUtil struct {
	validator *validator.Validate
}

// NewValidationUtil creates a new instance of ValidationUtil
func NewValidationUtil() *ValidationUtil {
	return &ValidationUtil{
		validator: validator.New(),
	}
}

// Validate performs struct validation and returns formatted HTTP errors
func (v *ValidationUtil) Validate(data interface{}) error {
	httpErrors := []errorHandler.HttpError{}

	err := v.validator.Struct(data)
	if err != nil {
		//nolint:errorlint
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return errorHandler.NewInternalServer()
		}

		for _, e := range validationErrors {
			httpError := errorHandler.HttpError{
				Code:    400,
				Info:    "Validation Error",
				Message: helper.FormatErrorMessage(e),
			}
			httpErrors = append(httpErrors, httpError)
		}
	}

	if len(httpErrors) > 0 {
		return errorHandler.NewMultipleFieldsValidation(httpErrors)
	}

	return nil
}
