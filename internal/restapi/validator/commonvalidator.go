package commonvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"auth-service-rizkysr90-pos/internal/constant"

	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"
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

func ValidatePassword(password string) *restapierror.RestAPIError {
	/*
	 that checks the length, presence of uppercase and lowercase letters,
	 digits, and special characters in the given password.
	 The contains function is used to check if the password contains characters
	 that satisfy a specific condition. Adjust the criteria as needed for your specific requirements.
	*/
	err := restapierror.RestAPIError{
		Code:    400,
		Message: "invalid password",
		//nolint:lll
		Details: "Password must be presence of uppercase and lowercase letters, digits, and special characters, min 8 char and max 64 char",
	}
	// Check length
	if len(password) < 8 {
		return &err
	}
	if len(password) > 64 {
		return &err
	}
	// Check for uppercase letter
	if !contains(password, isUpperCase) {
		return &err
	}

	// Check for lowercase letter
	if !contains(password, isLowerCase) {
		return &err
	}

	// Check for digit
	if !contains(password, isDigit) {
		return &err
	}

	// Check for special character
	if !contains(password, isSpecialChar) {
		return &err
	}

	return nil
}
func isUpperCase(r rune) bool {
	return unicode.IsUpper(r)
}

func isLowerCase(r rune) bool {
	return unicode.IsLower(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isSpecialChar(r rune) bool {
	return regexp.MustCompile(`[[:punct:]]`).MatchString(string(r))
}
func contains(s string, condition func(rune) bool) bool {
	for _, r := range s {
		if condition(r) {
			return true
		}
	}
	return false
}
