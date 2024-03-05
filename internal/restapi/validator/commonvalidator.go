package commonvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/rizkysr90/go-boilerplate/internal/constant"
	"github.com/rizkysr90/go-boilerplate/pkg/restapierror"
)

func ValidateRequired(value interface{}, field string) *restapierror.RestAPIError {
	err := restapierror.RestAPIError{
		Code:    400,
		Message: "REQUIRED DATA",
		Details: fmt.Sprintf("%s is required", field),
	}
	if value == nil {
		return &err
	}
	v := reflect.ValueOf(value)
	kind := reflect.TypeOf(value).Kind()

	// Getting underlying value for pointer
	if kind == reflect.Ptr {
		v = v.Elem()
		kind = v.Kind()
	}
	//nolint:exhaustive
	switch kind {
	case reflect.Array, reflect.Chan, reflect.Slice, reflect.Map:
		if v.Len() == 0 || v.IsZero() {
			return &err
		}
	case reflect.String:
		realValue := strings.TrimSpace(v.String())
		if len(realValue) == 0 {
			return &err
		}
	default:
		if v.IsZero() {
			return &err
		}
	}
	return nil
}
func ValidateEmail(email string, field string) *restapierror.RestAPIError {
	// Define a regular expression for basic email validation
	// This regex allows letters, numbers, dots, and underscores in the username part,
	// a single '@' symbol, and letters and dots in the domain part.
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// Use the MatchString method to check if the email matches the regex
	if !emailRegex.MatchString(email) {
		return &restapierror.RestAPIError{
			Code:    400,
			Message: constant.ErrInvalidFormat,
			Details: fmt.Sprintf("%s is required", field),
		}
	}
	return nil
}
