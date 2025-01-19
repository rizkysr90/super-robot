package utility

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plain text password and returns the bcrypt hash
func HashPassword(password string) (string, error) {
	// Generate hash with default cost of 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a plain text password against a hashed password
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
