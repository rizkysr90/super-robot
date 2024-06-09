package middleware

import (
	"auth-service-rizkysr90-pos/internal/helper/errorHandler"
	"bytes"

	"github.com/gin-gonic/gin"
)

// Custom response writer to capture the response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	// Write to the buffer and the original response writer
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func ResponseBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a custom response writer to capture the response body
		writer := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		// Proceed to the next middleware
		c.Next()
		// After the response is written, log the captured response body
		c.Set("response_body", writer.body.Bytes())

	}
}
func GetResBodyValue(ctx *gin.Context) ([]byte, error) {
	// Get the request body from the context
	resBody, ok := ctx.Get("response_body")
	if !ok {
		// Request body not found in context
		getResBodyErr := errorHandler.NewInternalServer(errorHandler.WithInfo("Response body not found in context"))
		return nil, getResBodyErr
	}
	return resBody.([]byte), nil
}
