package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rizkysr90-pos/pkg/errorHandler"

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
	//nolint:exhaustive
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
	default:
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

// ErrorHandler returns a middleware that handles errors and logs them appropriately.
// It uses the provided zerolog.Logger for logging.
func ErrorHandler(logger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Process request
		ctx.Next()

		// Skip if no errors
		if len(ctx.Errors) == 0 {
			return
		}

		// Get first error from context
		ginErr := ctx.Errors[0]

		var restAPIErr *errorHandler.HttpError
		if errors.As(ginErr.Err, &restAPIErr) {
			handleKnownError(ctx, logger, restAPIErr)
			return
		}

		handleUnknownError(ctx, logger, ginErr)
	}
}

// handleKnownError processes errors of type HttpError.
func handleKnownError(ctx *gin.Context, logger zerolog.Logger, err *errorHandler.HttpError) {
	loggerErrorFormat(ctx, logger, err, zerolog.WarnLevel)
	markErrorAsLogged(ctx)
	ctx.AbortWithStatusJSON(err.Code, err)
}

// handleUnknownError processes unknown error types.
func handleUnknownError(ctx *gin.Context, logger zerolog.Logger, ginErr *gin.Error) {
	httpErr := &errorHandler.HttpError{
		Code:    http.StatusInternalServerError,
		Info:    "Internal Server Error",
		Message: ginErr.Error(),
	}

	loggerErrorFormat(ctx, logger, httpErr, zerolog.ErrorLevel)
	markErrorAsLogged(ctx)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, httpErr)
}

// markErrorAsLogged marks the request as having its error logged.
func markErrorAsLogged(ctx *gin.Context) {
	ctx.Set("is_error_logged", true)
}
