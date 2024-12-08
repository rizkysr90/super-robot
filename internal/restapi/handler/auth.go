package handler

import (
	"net/http"
	"rizkysr90-pos/internal/auth"
	service "rizkysr90-pos/internal/service/auth"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authClient  *auth.Client
	authService *service.Auth
}

func NewAuthHandler(
	authClient *auth.Client,
	authService *service.Auth,
) *AuthHandler {
	return &AuthHandler{
		authClient:  authClient,
		authService: authService,
	}
}

func (a *AuthHandler) OwnerRegistration(ctx *gin.Context) {
	payload := &service.RequestRegisterOwner{}
	if err := ctx.ShouldBind(payload); err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}
	stateID, err := a.authService.RegisterOwner(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Redirect(http.StatusFound, a.authClient.Oauth.AuthCodeURL(stateID))
}
