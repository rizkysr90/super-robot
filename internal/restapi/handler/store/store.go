package store

import (
	"html"
	"net/http"
	"strings"

	payload "auth-service-rizkysr90-pos/internal/payload/http/store"
	commonvalidator "auth-service-rizkysr90-pos/internal/restapi/validator"
	"auth-service-rizkysr90-pos/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"
)

type StoreHandler struct {
	storeService service.StoreService
}

func NewAuthHandler(
	storeService service.StoreService,
) *StoreHandler {
	return &StoreHandler{
		storeService: storeService,
	}
}

type reqCreateStore struct {
	payload *payload.ReqCreateStore
}

func (req *reqCreateStore) sanitize() {
	req.payload.Name = html.EscapeString(strings.TrimSpace(req.payload.Name))
	req.payload.Address = html.EscapeString(strings.TrimSpace(req.payload.Address))
	req.payload.Contact = html.EscapeString(strings.TrimSpace(req.payload.Contact))
	req.payload.UserID = html.EscapeString(strings.TrimSpace(req.payload.UserID))
}
func (req *reqCreateStore) validate() error {
	validationErrors := []restapierror.RestAPIError{}
	if err := commonvalidator.ValidateRequired(req.payload.Name, "name"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.Address, "address"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.Contact, "contacts"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.UserID, "user_id"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if len(validationErrors) > 0 {
		return restapierror.NewMultipleFieldsValidation(validationErrors)
	}
	return nil
}
func (a *StoreHandler) CreateStore(ctx *gin.Context) {
	payload := &payload.ReqCreateStore{}
	payload.UserID = ctx.GetString("user_id")
	if err := ctx.Bind(payload); err != nil {
		ctx.Error(err)
		return
	}
	input := reqCreateStore{payload}
	input.sanitize()
	if err := input.validate(); err != nil {
		ctx.Error(err)
		return
	}
	if err := a.storeService.CreateStore(ctx, payload); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{})
}
