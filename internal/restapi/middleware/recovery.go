package middleware

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
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
				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) //nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, "error")
				}
				// Log the error
				// log.Println().Err(err.(error)).Msg("Middleware panic")
				// log.Println("panic err", err)
				// Respond with an internal server error
				// c.AbortWithStatusJSON(500, gin.H{
				// 	"message": "Internal Server Error",
				// })
			}
		}()
		// Call the next handler
		c.Next()
	}
}
