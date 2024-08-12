package validator

import (
	"auth-service-rizkysr90-pos/internal/constant"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

func ValidateName(name string) bool {
	// Regular expression for validating names
	// This regex allows alphabetic characters and optionally allows spaces, hyphens, and apostrophes
	nameRegex := regexp.MustCompile(`^[a-zA-Z]+([' -][a-zA-Z]+)*$`)
	return nameRegex.MatchString(name)
}
func ValidateRequired(value interface{}, field string) bool {
	if value == nil {
		return false
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
			return false
		}
	case reflect.String:
		realValue := strings.TrimSpace(v.String())
		if len(realValue) == 0 {
			return false
		}
	default:
		if v.IsZero() {
			return false
		}
	}
	return true
}
func ValidateEmail(email string) bool {
	// Define a regular expression for basic email validation
	// This regex allows letters, numbers, dots, and underscores in the username part,
	// a single '@' symbol, and letters and dots in the domain part.
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// Use the MatchString method to check if the email matches the regex
	return emailRegex.MatchString(email)
}
func ValidateOnlyNumber(input string) bool {
	// Regular expression to match only digits (0-9)
	regex := regexp.MustCompile(`^[0-9]+$`)
	return regex.MatchString(input)
}
func ValidateRoles(role int) bool {
	values := []int{constant.RBAC_ADMIN, constant.RBAC_OPERATIONAL}
	for _, v := range values {
		if role == v {
			return true // Input value found in array, so it's valid
		}
	}
	return false // Input value not found in array, so it's invalid
}
func ValidatePassword(password string) bool {
	/*
	 that checks the length, presence of uppercase and lowercase letters,
	 digits, and special characters in the given password.
	 The contains function is used to check if the password contains characters
	 that satisfy a specific condition. Adjust the criteria as needed for your specific requirements.
	*/

	// Check length
	if len(password) < 8 {
		return false
	}
	if len(password) > 64 {
		return false
	}
	// Check for uppercase letter
	if !contains(password, isUpperCase) {
		return false
	}

	// Check for lowercase letter
	if !contains(password, isLowerCase) {
		return false
	}

	// Check for digit
	if !contains(password, isDigit) {
		return false
	}

	// Check for special character
	if !contains(password, isSpecialChar) {
		return false
	}

	return true
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
