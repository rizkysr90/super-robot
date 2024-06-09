package utility

import (
	"strings"

	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"
)

func ReturnAPIError(code int, message, details string) *restapierror.RestAPIError {
	return &restapierror.RestAPIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
func SanitizeReqBody(input string) string {
	sanitizedReqBody := strings.ReplaceAll(input, "\r", "")
	sanitizedReqBody = strings.ReplaceAll(sanitizedReqBody, "\n", "")
	sanitizedReqBody = strings.ReplaceAll(sanitizedReqBody, "\\", "")
	sanitizedReqBody = strings.ReplaceAll(sanitizedReqBody, " ", "")
	return sanitizedReqBody
}
