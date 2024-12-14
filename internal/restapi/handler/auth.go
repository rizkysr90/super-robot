package handler

import (
	"log"
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

// OwnerRegistration godoc
// @Summary Register owner for a tenant
// @Description Initiates owner registration process and redirects to OAuth authorization
// @Tags auth
// @Accept json
// @Produce json
// @Param tenant query string true "Tenant name to register owner for"
// @Success 302 {string} string "Redirect to OAuth authorization URL"
// @Failure 400 {object} errorHandler.HttpError
// @Failure 500 {object} errorHandler.HttpError
// @Router /auth/register/owner [get]
func (a *AuthHandler) OwnerRegistration(ctx *gin.Context) {
	tenantName := ctx.Query("tenant")
	if tenantName == "" {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage("tenantName is required"),
		)
		if errCtx := ctx.Error(err); errCtx != nil {
			// Handle the error from ctx.Error
			// You might want to log it or take other appropriate action
			log.Printf("failed to send error response: %v", errCtx)
		}
		return
	}
	payload := &service.RequestRegisterOwner{TenantName: tenantName}
	stateID, err := a.authService.RegisterOwner(ctx, payload)
	if err != nil {
		if errCtx := ctx.Error(err); errCtx != nil {
			// Handle the error from ctx.Error
			// You might want to log it or take other appropriate action
			log.Printf("failed to send error response: %v", errCtx)
		}
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
		if errCtx := ctx.Error(err); errCtx != nil {
			// Handle the error from ctx.Error
			// You might want to log it or take other appropriate action
			log.Printf("failed to send error response: %v", errCtx)
		}
		return
	}
	authorizationCode := ctx.Query("code")
	if authorizationCode == "" {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid code"),
			errorHandler.WithMessage("authorization code is required"),
		)
		if errCtx := ctx.Error(err); errCtx != nil {
			// Handle the error from ctx.Error
			// You might want to log it or take other appropriate action
			log.Printf("failed to send error response: %v", errCtx)
		}
		return
	}
	response, err := a.authService.Callback(ctx, &service.RequestCallback{
		State: stateID, Code: authorizationCode})
	if err != nil {
		if errCtx := ctx.Error(err); errCtx != nil {
			// Handle the error from ctx.Error
			// You might want to log it or take other appropriate action
			log.Printf("failed to send error response: %v", errCtx)
		}
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
