package employee

import (
	"html"
	"net/http"
	"strconv"
	"strings"

	payload "auth-service-rizkysr90-pos/internal/payload/http/employee"
	commonvalidator "auth-service-rizkysr90-pos/internal/restapi/validator"
	"auth-service-rizkysr90-pos/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"
)

type EmployeeHandler struct {
	employeeService service.EmployeeService
}

func NewEmployeeHandler(
	employeeService service.EmployeeService,
) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

type reqCreateEmployee struct {
	payload *payload.ReqCreateEmployee
}

func (req *reqCreateEmployee) sanitize() {
	req.payload.Name = html.EscapeString(strings.TrimSpace(req.payload.Name))
	req.payload.Contact = html.EscapeString(strings.TrimSpace(req.payload.Contact))
	req.payload.Username = html.EscapeString(strings.TrimSpace(req.payload.Username))
	req.payload.Password = html.EscapeString(strings.TrimSpace(req.payload.Password))
	req.payload.ConfirmPassword = html.EscapeString(strings.TrimSpace(req.payload.ConfirmPassword))
	req.payload.StoreID = html.EscapeString(strings.TrimSpace(req.payload.StoreID))

}
func (req *reqCreateEmployee) validate() error {
	validationErrors := []restapierror.RestAPIError{}
	if err := commonvalidator.ValidateRequired(req.payload.Name, "name"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.Contact, "contact"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.Username, "username"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.Password, "password"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.ConfirmPassword, "confirm_password"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.StoreID, "store_id"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.payload.Role, "role"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if len(validationErrors) > 0 {
		return restapierror.NewMultipleFieldsValidation(validationErrors)
	}
	return nil
}
func (e *EmployeeHandler) CreateStore(ctx *gin.Context) {
	payload := &payload.ReqCreateEmployee{}
	if err := ctx.Bind(payload); err != nil {
		ctx.Error(err)
		return
	}
	input := reqCreateEmployee{payload}
	input.sanitize()
	if err := input.validate(); err != nil {
		ctx.Error(err)
		return
	}
	if err := e.employeeService.Create(ctx, payload); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{})
}

type reqLoginEmployee struct {
	*payload.ReqLoginEmployee
}

func (req *reqLoginEmployee) sanitize() {
	req.Username = html.EscapeString(strings.TrimSpace(req.Username))
	req.Password = html.EscapeString(strings.TrimSpace(req.Password))
}
func (req *reqLoginEmployee) validate() error {
	validationErrors := []restapierror.RestAPIError{}
	if err := commonvalidator.ValidateRequired(req.Username, "username"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := commonvalidator.ValidateRequired(req.Password, "password"); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if len(validationErrors) > 0 {
		return restapierror.NewMultipleFieldsValidation(validationErrors)
	}
	return nil
}
func (e *EmployeeHandler) LoginUser(ctx *gin.Context) {
	payload := &payload.ReqLoginEmployee{}
	if err := ctx.ShouldBindJSON(payload); err != nil {
		ctx.Error(err)
		return
	}
	input := reqLoginEmployee{payload}
	input.sanitize()
	if err := input.validate(); err != nil {
		ctx.Error(err)
		return
	}
	data, err := e.employeeService.Login(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	getRole := strconv.Itoa(data.Role)
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("refresh_token", data.RefreshToken, 0, "", "", true, true)
	ctx.SetCookie("user_role", getRole, 0, "", "", true, false)
	ctx.SetCookie("access_token", data.Token, 0, "", "", true, true)
	ctx.JSON(http.StatusOK, gin.H{})

}

// func (a *StoreHandler) GetAllStore(ctx *gin.Context) {
// 	var err error
// 	req := &payload.ReqGetAllStore{
// 		UserID: ctx.GetString("user_id"),
// 	}
// 	req.PageNumber, err = strconv.Atoi(ctx.Query("page_number"))
// 	if err != nil {
// 		req.PageNumber = 0
// 	}
// 	req.PageSize, err = strconv.Atoi(ctx.Query("page_size"))
// 	if err != nil {
// 		req.PageSize = 0
// 	}
// 	var responseData *payload.ResGetAllStore
// 	if responseData, err = a.storeService.GetAllStore(ctx, req); err != nil {
// 		ctx.Error(err)
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"pagination": responseData.Pagination,
// 		"data":       responseData.Data,
// 	})
// }
