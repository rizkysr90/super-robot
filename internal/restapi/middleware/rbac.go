package middleware

import (
	"net/http"

	"auth-service-rizkysr90-pos/pkg/errorHandler"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware is a middleware function to enforce RBAC
func RBACMiddleware(userRoleLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Assuming user roles are stored in the context after authentication
		userRoles := c.GetInt("userRoles")
		if userRoles == 0 {
			err := errorHandler.NewUnauthorized(errorHandler.WithMessage("User roles not found"))
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		// Check if user has any of the allowed roles
		if userRoles <= userRoleLevel {
			c.Next()
		} else {
			err := errorHandler.NewUnauthorized(errorHandler.WithMessage("not allowed"))
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
	}
}
