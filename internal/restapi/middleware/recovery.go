package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rizkysr90-pos/pkg/errorHandler"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Recovery(logger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				handlePanic(ctx, logger, err)
			}
		}()
		ctx.Next()
	}
}

// handlePanic processes recovered panic and logs error details
func handlePanic(ctx *gin.Context, logger zerolog.Logger, recoveredErr interface{}) {
	if err, ok := recoveredErr.(error); ok {
		httpErr := createHTTPError(err)
		logPanicDetails(ctx, logger, httpErr)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, httpErr)
	}
}

// createHttpError wraps the error in a standardized HTTP error response
func createHTTPError(err error) errorHandler.HttpError {
	return errorHandler.HttpError{
		Code:    http.StatusInternalServerError,
		Info:    "Panic - Internal Server Error",
		Message: err.Error(),
	}
}

// logPanicDetails logs all relevant information about the panic
func logPanicDetails(ctx *gin.Context, logger zerolog.Logger, httpErr errorHandler.HttpError) {
	reqBody := getRequestBody(ctx)
	resBody := getResponseBody(httpErr)
	stackTrace := debug.Stack()

	logger.Error().
		Str("method", ctx.Request.Method).
		Str("path", ctx.Request.URL.String()).
		Str("client_ip", ctx.ClientIP()).
		Str("user_agent", ctx.GetHeader("User-Agent")).
		Int("statusCode", httpErr.Code).
		RawJSON("req_body", reqBody).
		RawJSON("res_body", resBody).
		Str("stack", string(stackTrace)).
		Msg(fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.String()))
}

// getRequestBody safely retrieves the request body or returns empty bytes
func getRequestBody(ctx *gin.Context) []byte {
	body, err := GetRequstBodyValue(ctx)
	if err != nil {
		return []byte{}
	}
	return body
}

// getResponseBody safely converts HTTP error to JSON or returns empty bytes
func getResponseBody(httpErr errorHandler.HttpError) []byte {
	body, err := json.Marshal(httpErr)
	if err != nil {
		return []byte{}
	}
	return body
}
