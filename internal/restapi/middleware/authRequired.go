package middleware

import (
	"strings"

	"github.com/rizkysr90/go-boilerplate/internal/config"
	jwttoken "github.com/rizkysr90/go-boilerplate/pkg/jwt"
	"github.com/rizkysr90/go-boilerplate/pkg/restapierror"

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
