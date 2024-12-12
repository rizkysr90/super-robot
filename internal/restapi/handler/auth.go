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
	tenantName := ctx.Query("tenant")
	if tenantName == "" {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage("tenantName is required"),
		)
		ctx.Error(err)
		return
	}
	payload := &service.RequestRegisterOwner{TenantName: tenantName}
	stateID, err := a.authService.RegisterOwner(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Redirect(http.StatusFound, a.authClient.Oauth.AuthCodeURL(stateID))
}
func (a *AuthHandler) Callback(ctx *gin.Context) {
	stateID := ctx.Query("state")
	if stateID == "" {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage("state id is required"),
		)
		ctx.Error(err)
		return
	}
	authorizationCode := ctx.Query("code")
	if authorizationCode == "" {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid code"),
			errorHandler.WithMessage("authorization code is required"),
		)
		ctx.Error(err)
		return
	}
	response, err := a.authService.Callback(ctx, &service.RequestCallback{
		State: stateID, Code: authorizationCode})
	if err != nil {
		ctx.Error(err)
		return
	}

	// Set session cookie
	// Parameters:
	// 1. name: cookie name
	// 2. value: cookie value (userID/session ID)
	// 3. maxAge: cookie duration in seconds
	// 4. path: cookie path
	// 5. domain: cookie domain
	// 6. secure: only send cookie over HTTPS
	// 7. httpOnly: prevent JavaScript access to cookie
	ctx.SetCookie(
		"session_id",                   // name
		response.SessionData.SessionID, // value
		3600,                           // maxAge (1 hour)
		"/",                            // path
		"",                             // domain
		true,                           // secure
		true,                           // httpOnly
	)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")

}
