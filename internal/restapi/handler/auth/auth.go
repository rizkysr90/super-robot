package auth

import (
	"html"
	"net/http"
	"strings"

	"github.com/rizkysr90/go-boilerplate/internal/config"
	payload "github.com/rizkysr90/go-boilerplate/internal/payload/http/auth"
	commonvalidator "github.com/rizkysr90/go-boilerplate/internal/restapi/validator"
	"github.com/rizkysr90/go-boilerplate/internal/service"
	"github.com/rizkysr90/go-boilerplate/pkg/restapierror"

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
func (req *reqCreateUser) validate() error {
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
	if err := commonvalidator.ValidateRequired(req.ConfirmPassword, "confirm_password"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateEmail(req.Email, "email"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidatePassword(req.Password); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if len(validationErrors) > 0 {
		return restapierror.NewMultipleFieldsValidation(validationErrors)
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
	if err := input.validate(); err != nil {
		ctx.Error(err)
		return
	}
	if err := a.authService.CreateUser(ctx, payload); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{})
}

type reqLoginUser struct {
	*payload.ReqLoginUser
}

func (req *reqLoginUser) sanitize() {
	req.Email = html.EscapeString(strings.TrimSpace(req.Email))
	req.Password = html.EscapeString(strings.TrimSpace(req.Password))
}
func (req *reqLoginUser) validate() error {
	validationErrors := []restapierror.RestAPIError{}
	if err := commonvalidator.ValidateRequired(req.Email, "email"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.Password, "password"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateEmail(req.Email, "email"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if len(validationErrors) > 0 {
		return restapierror.NewMultipleFieldsValidation(validationErrors)
	}
	return nil
}
func (a *AuthHandler) LoginUser(ctx *gin.Context) {
	payload := &payload.ReqLoginUser{}
	if err := ctx.ShouldBindJSON(payload); err != nil {
		ctx.Error(err)
		return
	}
	input := reqLoginUser{payload}
	input.sanitize()
	if err := input.validate(); err != nil {
		ctx.Error(err)
		return
	}
	data, err := a.authService.LoginUser(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)

}
