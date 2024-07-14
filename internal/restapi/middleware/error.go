package middleware

import (
	"auth-service-rizkysr90-pos/internal/helper/errorHandler"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func loggerErrorFormat(ctx *gin.Context,
	logger zerolog.Logger,
	httpErr *errorHandler.HttpError,
	level zerolog.Level) {
	getReqBody, errReqBody := GetRequstBodyValue(ctx)
	if errReqBody != nil {
		panic(errReqBody)
	}
	// Convert the struct to JSON format
	jsonResBody, errResBody := json.Marshal(httpErr)
	if errResBody != nil {
		panic(errResBody)
	}
	// sanitizedReqBody := utility.SanitizeReqBody(getReqBody)
	switch level {
	case zerolog.WarnLevel:
		logger.Warn().
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.String()).
			Str("client_ip", ctx.ClientIP()).
			Str("user_agent", ctx.GetHeader("User-Agent")).
			Int("statusCode", httpErr.Code).
			RawJSON("req_body", getReqBody).
			RawJSON("res_body", jsonResBody).
			Msg(fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.String()))
	case zerolog.ErrorLevel:

		// log.Printf("Stack Trace:\n%s\n", string(stackTrace))
		logger.Error().
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.String()).
			Str("client_ip", ctx.ClientIP()).
			Str("user_agent", ctx.GetHeader("User-Agent")).
			Int("statusCode", httpErr.Code).
			RawJSON("req_body", getReqBody).
			RawJSON("res_body", jsonResBody).
			Msg(fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.String()))
	}

}
func ErrorHandler(logger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		// Check if there are errors
		hasErrors := len(ctx.Errors) > 0
		if hasErrors {
			ginErr := ctx.Errors[0]

			// Type assertion from gin.Error to errorHandler.HttpError
			if restAPIErr, ok := ginErr.Err.(*errorHandler.HttpError); ok {
				loggerErrorFormat(ctx, logger, restAPIErr, zerolog.WarnLevel)
				// Mark request as processed successfully
				ctx.Set("is_error_logged", bool(true))
				ctx.AbortWithStatusJSON(restAPIErr.Code, restAPIErr)

			} else {
				httpErr := errorHandler.HttpError{
					Code:    500,
					Info:    "Internal Server Error",
					Message: ginErr.Error(),
				}
				loggerErrorFormat(ctx, logger, &httpErr, zerolog.ErrorLevel)
				// Mark request as processed successfully
				ctx.Set("is_error_logged", bool(true))
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, httpErr)

			}
			return
		}
	}
}

// func GetIsErrorLogged(ctx *gin.Context) (bool, error) {
// 	// Get the request body from the context
// 	value := ctx.GetBool("is_error_logged")
// 	if !ok {
// 		// Request body not found in context
// 		getRequestBodyErr := errorHandler.NewInternalServer(
// 			errorHandler.WithInfo("is_error_logged not found in context"))
// 		return nil, getRequestBodyErr
// 	}
// 	return requestBody.([]byte), nil
// }
