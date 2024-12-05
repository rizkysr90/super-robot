package utility

import (
	"crypto/rand"
	"encoding/base64"
	"rizkysr90-pos/pkg/errorHandler"
	"strings"
)

func ReturnAPIError(code int, message, info string) *errorHandler.HttpError {
	return &errorHandler.HttpError{
		Code:    code,
		Message: message,
		Info:    info,
	}
}
func SanitizeReqBody(input string) string {
	sanitizedReqBody := strings.ReplaceAll(input, "\r", "")
	sanitizedReqBody = strings.ReplaceAll(sanitizedReqBody, "\n", "")
	sanitizedReqBody = strings.ReplaceAll(sanitizedReqBody, "\\", "")
	sanitizedReqBody = strings.ReplaceAll(sanitizedReqBody, " ", "")
	return sanitizedReqBody
}
func GenerateRandomBase64Str() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
