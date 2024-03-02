package auth

import (
	"api-iad-ams/internal/config"
	payload "api-iad-ams/internal/payload/http/auth"
	commonvalidator "api-iad-ams/internal/restapi/validator"
	"api-iad-ams/internal/service"
	"api-iad-ams/pkg/restapierror"
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	// db          *pgxpool.Pool
	authService service.AuthService
	config      config.Config
}

func NewAuthHandler(
	authService service.AuthService,
	config config.Config,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}
func (a *AuthHandler) AddRoutes(ginEngine *gin.Engine) {
	createUserPath := fmt.Sprintf("%s/auth/users", a.config.ApiVersionBaseURL)
	ginEngine.POST(createUserPath, func(ctx *gin.Context) {
		a.CreateUser(ctx)
	})
}

//	type testError struct {
//		TestArray string `json:"test_array"`
//	}
type reqCreateUser struct {
	*payload.ReqCreateAccount
}

func (req *reqCreateUser) sanitize() {
	req.FirstName = html.EscapeString(strings.TrimSpace(req.FirstName))
	req.LastName = html.EscapeString(strings.TrimSpace(req.LastName))
	req.Email = html.EscapeString(strings.TrimSpace(req.Email))
	req.Password = html.EscapeString(strings.TrimSpace(req.Password))
	req.ConfirmPassword = html.EscapeString(strings.TrimSpace(req.ConfirmPassword))
}
func (req *reqCreateUser) validate(ctx context.Context) error {
	validationErrors := []restapierror.RestAPIError{}
	if err := commonvalidator.ValidateRequired(req.FirstName, "first_name"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.LastName, "last_name"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.Email, "email"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.Password, "password"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.ConfirmPassword, "confirm_password"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateEmail(req.Email, "email"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if len(validationErrors) > 0 {
		return restapierror.NewMultipleFieldsValidation(ctx, validationErrors)
	}

	return nil
}
func (a *AuthHandler) CreateUser(ctx *gin.Context) {
	payload := &payload.ReqCreateAccount{}

	if err := ctx.Bind(payload); err != nil {
		ctx.Error(err)
		return
	}
	input := reqCreateUser{payload}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		ctx.Error(err)
		return
	}
	if err := a.authService.CreateUser(ctx, payload); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{})
}
