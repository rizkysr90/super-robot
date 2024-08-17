package utility

import (
	"auth-service-rizkysr90-pos/pkg/errorHandler"
	"strings"
)

func ReturnAPIError(code int, message, info string) *errorHandler.HttpError {
	return &errorHandler.HttpError{
		Code:    code,
		Message: message,
		Info: info,
	}
}
func SanitizeReqBody(input string) string {
	sanitizedReqBody := strings.ReplaceAll(input, "\r", "")
	sanitizedReqBody = strings.ReplaceAll(sanitizedReqBody, "\n", "")
	sanitizedReqBody = strings.ReplaceAll(sanitizedReqBody, "\\", "")
	sanitizedReqBody = strings.ReplaceAll(sanitizedReqBody, " ", "")
	return sanitizedReqBody
}
