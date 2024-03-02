package middleware

import (
	"api-iad-ams/pkg/restapierror"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func logRestAPIErr(c *gin.Context, logger zerolog.Logger, restAPIErr *restapierror.RestAPIError) {
	logger.Error().
		Str("code", fmt.Sprintf("%d", restAPIErr.Code)).
		Str("path", c.FullPath()).
		Any("details", restAPIErr.Details).
		Msg(restAPIErr.Message)
}
func ErrorHandler(logger zerolog.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()
		// Check if there are errors
		hasErrors := len(c.Errors) > 0
		if hasErrors {
			ginErr := c.Errors[0]
			if restAPIErr, ok := ginErr.Err.(*restapierror.RestAPIError); ok {
				c.AbortWithStatusJSON(restAPIErr.Code, restAPIErr)
				logRestAPIErr(c, logger, restAPIErr)
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Internal Server Error",
				"details": ginErr.Err.Error(),
			})
			logger.Error().
				Int("code", http.StatusInternalServerError).
				Str("path", c.FullPath()).
				Msg(ginErr.Error())
		}
		if !hasErrors {
			logger.Info().
				Str("path", c.FullPath()).
				Msg("Request processed")
		}

	}
}
