package middleware

import (
	"errors"
	"net/http"

	"github.com/rizkysr90/go-boilerplate/internal/config"
	"github.com/rizkysr90/go-boilerplate/internal/constant"
	"github.com/rizkysr90/go-boilerplate/internal/helper"

	"github.com/gin-gonic/gin"
)

func AuthRequired(cfg config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getAPIKey := ctx.Request.Header.Get("API_KEY")
		if getAPIKey != cfg.APIKey {
			err := errors.New("contact administrator to get API_KEY")
			ctx.Error(err)
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				helper.RestApiError(http.StatusUnauthorized, constant.ERR_INVALID_API_KEY, err.Error()),
			)
			return
		}
		ctx.Next()
	}
}
