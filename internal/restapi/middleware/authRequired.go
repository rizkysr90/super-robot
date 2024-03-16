package middleware

import (
	"strings"

	"auth-service-rizkysr90-pos/internal/config"
	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"

	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"

	"github.com/gin-gonic/gin"
)

func AuthRequired(cfg config.Config, jwt *jwttoken.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getAuthHeader := ctx.Request.Header.Get("Authorization")
		splitAuthHeader := strings.Split(getAuthHeader, " ")
		if len(splitAuthHeader) < 2 {
			err := restapierror.NewUnauthorized(restapierror.WithMessage("access token not provided"))
			ctx.Error(err)
			return
		}
		getToken := splitAuthHeader[1]
		if err := jwt.Authorize(getToken); err != nil {
			ctx.Error(err)
			return
		}
		ctx.Next()
	}
}
