package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func LogMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		getReqBody, errReqBody := GetRequstBodyValue(ctx)
		if errReqBody != nil {
			panic(errReqBody)
		}
		// Convert the struct to JSON format
		getResBody, errResBody := GetResBodyValue(ctx)
		if errResBody != nil {
			panic(errResBody)
		}
		if !ctx.GetBool("is_error_logged") {
			logger.Info().
				Str("method", ctx.Request.Method).
				Str("path", ctx.Request.URL.String()).
				Str("client_ip", ctx.ClientIP()).
				Str("user_agent", ctx.GetHeader("User-Agent")).
				Int("statusCode", ctx.Writer.Status()).
				RawJSON("req_body", getReqBody).
				RawJSON("res_body", getResBody).
				Msg(fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.String()))
		}
	}
}
