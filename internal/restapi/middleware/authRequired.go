package middleware

import (
	"net/http"

	"github.com/rizkysr90/go-boilerplate/internal/config"
	"github.com/rizkysr90/go-boilerplate/pkg/restapierror"

	"github.com/gin-gonic/gin"
)

func AuthRequired(cfg config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getAPIKey := ctx.Request.Header.Get("API_KEY")
		if getAPIKey != cfg.APIKey {
			err := restapierror.NewUnauthorized(restapierror.WithMessage("contact administrator to get API_KEY"))
			if errCtx := ctx.Error(err); errCtx != nil {
				// Handle the error. This could mean logging it, returning an error response to the client,
				// or even just panicking, depending on the severity of the error and the design of your application.
				// Optionally, you can also send a response to the client indicating an internal server error.
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}
		}
		ctx.Next()
	}
}
