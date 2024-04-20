package store

import (
	"html"
	"net/http"
	"strconv"
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
	req.payload.EmployeeID = html.EscapeString(strings.TrimSpace(req.payload.EmployeeID))
}
func (req *reqCreateStore) validate() error {
	validationErrors := []restapierror.RestAPIError{}
	if err := commonvalidator.ValidateRequired(req.payload.Name, "name"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.Address, "address"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.Contact, "contact"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.EmployeeID, "user_id"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateOnlyNumber(req.payload.Contact, "contact"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if len(validationErrors) > 0 {
		return restapierror.NewMultipleFieldsValidation(validationErrors)
	}
	return nil
}
func (a *StoreHandler) CreateStore(ctx *gin.Context) {
	payload := &payload.ReqCreateStore{}
	if err := ctx.Bind(payload); err != nil {
		ctx.Error(err)
		return
	}
	payload.EmployeeID = ctx.GetString("user_id")
	input := reqCreateStore{payload}
	input.sanitize()
	if err := input.validate(); err != nil {
		ctx.Error(err)
		return
	}
	data, err := a.storeService.CreateStore(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"store_id": data.StoreID, "created_at": data.CreatedAt})
}
func (a *StoreHandler) GetAllStore(ctx *gin.Context) {
	var err error
	req := &payload.ReqGetAllStore{
		EmployeeID: ctx.GetString("user_id"),
	}
	req.PageNumber, err = strconv.Atoi(ctx.Query("page_number"))
	if err != nil {
		req.PageNumber = 0
	}
	req.PageSize, err = strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		req.PageSize = 0
	}
	var responseData *payload.ResGetAllStore
	if responseData, err = a.storeService.GetAllStore(ctx, req); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"pagination": responseData.Pagination,
		"data":       responseData.Data,
	})
}

func (a *StoreHandler) DeleteStore(ctx *gin.Context) {
	payload := &payload.ReqDeleteStore{}
	if err := ctx.Bind(payload); err != nil {
		ctx.Error(err)
		return
	}
	validationErrors := []restapierror.RestAPIError{}
	payload.EmployeeID = ctx.GetString("user_id")
	payload.StoreID = strings.TrimSpace(payload.StoreID)
	payload.EmployeeID = strings.TrimSpace(payload.EmployeeID)
	if err := commonvalidator.ValidateRequired(payload.EmployeeID, "employee_id"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(payload.StoreID, "store_id"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if len(validationErrors) > 0 {
		ctx.Error(restapierror.NewMultipleFieldsValidation(validationErrors))
		return
	}
	if err := a.storeService.DeleteStore(ctx, payload); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

type reqUpdateStore struct {
	payload *payload.ReqUpdateStore
}

func (req *reqUpdateStore) sanitize() {
	req.payload.Name = html.EscapeString(strings.TrimSpace(req.payload.Name))
	req.payload.Address = html.EscapeString(strings.TrimSpace(req.payload.Address))
	req.payload.Contact = html.EscapeString(strings.TrimSpace(req.payload.Contact))
	req.payload.StoreID = html.EscapeString(strings.TrimSpace(req.payload.StoreID))
}
func (req *reqUpdateStore) validate() error {
	validationErrors := []restapierror.RestAPIError{}
	if err := commonvalidator.ValidateRequired(req.payload.Name, "name"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.Address, "address"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.Contact, "contact"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.StoreID, "store_id"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateOnlyNumber(req.payload.Contact, "contact"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if len(validationErrors) > 0 {
		return restapierror.NewMultipleFieldsValidation(validationErrors)
	}
	return nil
}
func (a *StoreHandler) UpdateStore(ctx *gin.Context) {
	payload := &payload.ReqUpdateStore{}
	if err := ctx.Bind(payload); err != nil {
		ctx.Error(err)
		return
	}
	// Get the value of store_id from the route path
	store_id := ctx.Param("store_id")
	payload.EmployeeID = store_id
	input := reqUpdateStore{payload}
	input.sanitize()
	if err := input.validate(); err != nil {
		ctx.Error(err)
		return
	}
	err := a.storeService.UpdateStore(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{})
}
