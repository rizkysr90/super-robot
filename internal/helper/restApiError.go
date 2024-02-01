package helper

import "github.com/gin-gonic/gin"

func RestApiError(code int, message string, details interface{}) *gin.H {
	return &gin.H{
		"code":    code,
		"message": message,
		"details": details,
	}
}
