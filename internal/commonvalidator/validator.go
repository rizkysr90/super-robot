// commonvalidator/validator.go
package commonvalidator

import (
	"regexp"
	"strings"
	"unicode"
)

// Regular expression for email validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsEmail checks if a string is a valid email address
func IsEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

// IsStrongPassword checks if a password meets minimum security requirements:
// - At least 8 characters long
// - Contains at least one uppercase letter
// - Contains at least one lowercase letter
// - Contains at least one number
// - Contains at least one special character
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// IsUserType checks if a string is a valid user type
func IsUserType(userType string) bool {
	validTypes := map[string]bool{
		"owner":  true,
		"admin":  true,
		"branch": true,
	}

	return validTypes[strings.ToLower(strings.TrimSpace(userType))]
}
