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
			//nolint:nestif
			if err := recover(); err != nil {
				if actualErr, ok := err.(error); ok {
					// Handle the error-of-error handling situation.
					httpErr := errorHandler.HttpError{
						Code:    500,
						Info:    "Panic - Internal Server Error",
						Message: actualErr.Error(),
					}
					getReqBody, errReqBody := GetRequstBodyValue(ctx)
					if errReqBody != nil {
						getReqBody = []byte{}
					}
					// Convert the struct to JSON format
					jsonResBody, errResBody := json.Marshal(httpErr)
					if errResBody != nil {
						jsonResBody = []byte{}
					}
					// Log the stack trace
					stackTrace := debug.Stack()
					// log.Printf("Stack Trace:\n%s\n", string(stackTrace))
					logger.Error().
						Str("method", ctx.Request.Method).
						Str("path", ctx.Request.URL.String()).
						Str("client_ip", ctx.ClientIP()).
						Str("user_agent", ctx.GetHeader("User-Agent")).
						Int("statusCode", httpErr.Code).
						RawJSON("req_body", getReqBody).
						RawJSON("res_body", jsonResBody).
						Str("stack", string(stackTrace)).
						Msg(fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.String()))
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, httpErr)
				}
			}
		}()
		// Call the next handler
		ctx.Next()
	}
}
