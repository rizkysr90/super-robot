package middleware

import (
	"strings"

	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"

	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"

	"github.com/gin-gonic/gin"
)

func AuthRequired(jwt *jwttoken.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getAuthHeader := ctx.Request.Header.Get("Authorization")
		splitAuthHeader := strings.Split(getAuthHeader, " ")
		if len(splitAuthHeader) < 2 {
			err := restapierror.NewUnauthorized(restapierror.WithMessage("access token not provided"))
			ctx.Error(err)
			return
		}
		getToken := splitAuthHeader[1]
		if _, err := jwt.Authorize(getToken); err != nil {
			ctx.Error(err)
			return
		}
		ctx.Next()
	}
}
func AuthRequiredCookies(jwt *jwttoken.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		access_token, err := ctx.Request.Cookie("access_token")
		if err != nil {
			ctx.Error(restapierror.NewUnauthorized(restapierror.WithMessage(err.Error())))
			return
		}
		if access_token.Value != "" {
			// Token found in cookies, authorize using it
			claims, err := jwt.Authorize(access_token.Value)
			if err != nil {
				// authorization error
				err = restapierror.NewUnauthorized(restapierror.WithMessage(err.Error()))
				ctx.Error(err)
				return
			}
			ctx.Set("user_id", claims.Subject)
			// Authorization successful, proceed with the request
			ctx.Next()
		} else {
			// Token not found in cookies or failed to authorize, check for token in headers
			token := ctx.GetHeader("Authorization")
			if token == "" {
				err := restapierror.NewUnauthorized(restapierror.WithMessage("access token not provided"))
				ctx.Error(err)
				return
			}
			// Remove "Bearer " prefix if present
			token = strings.TrimPrefix(token, "Bearer ")
			if _, err := jwt.Authorize(token); err != nil {
				ctx.Error(restapierror.NewUnauthorized(restapierror.WithMessage(err.Error())))
				return
			}
			ctx.Next()
		}
	}
}
