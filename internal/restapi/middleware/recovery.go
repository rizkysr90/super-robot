package middleware

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Recovery(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if !brokenPipe {
					c.Error(err.(error))
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"code":    500,
						"message": "Internal Server Error",
						"details": err.(error).Error(),
					})
				}
				logger.Error().
					Str("code", "500").
					Str("path", c.FullPath()).
					Msg(err.(error).Error())
			}
		}()
		// Call the next handler
		c.Next()
	}
}
