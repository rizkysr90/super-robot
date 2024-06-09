package middleware

import (
	"auth-service-rizkysr90-pos/internal/helper/errorHandler"
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestBodyMiddleware is a middleware function that reads the request body
// and stores it in the context.
func RequestBodyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the request body
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			// Handle error (e.g., logging)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}

		// Reset the request body so it can be read again
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		// Store the request body in the context
		c.Set("request_body", requestBody)

		// Call the next handler
		c.Next()
	}
}
func GetRequstBodyValue(ctx *gin.Context) ([]byte, error) {
	// Get the request body from the context
	requestBody, ok := ctx.Get("request_body")
	if !ok {
		// Request body not found in context
		getRequestBodyErr := errorHandler.NewInternalServer(errorHandler.WithInfo("Request body not found in context"))
		return nil, getRequestBodyErr
	}
	return requestBody.([]byte), nil
}
