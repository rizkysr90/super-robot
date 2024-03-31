package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"
)

// RBACMiddleware is a middleware function to enforce RBAC
func RBACMiddleware(userRoleLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Assuming user roles are stored in the context after authentication
		userRoles := c.GetInt("userRoles")
		if userRoles == 0 {
			err := restapierror.NewUnauthorized(restapierror.WithMessage("User roles not found"))
			c.Error(err)
			return
		}

		// Check if user has any of the allowed roles
		if userRoles <= userRoleLevel {
			c.Next()
		} else {
			err := restapierror.NewUnauthorized(restapierror.WithMessage("not allowed"))
			c.Error(err)
			return
		}
	}
}
