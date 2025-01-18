package handler

import (
	"net/http"
	"rizkysr90-pos/internal/service/admin"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/gin-gonic/gin"
)

type UserAdmin struct {
	adminService *admin.Service
}

func NewUserAdmin(adminService *admin.Service) *UserAdmin {
	return &UserAdmin{adminService: adminService}
}

// @Summary Create User
// @Description Creates a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body admin.RequestCreateAdmin true "User creation request"
// @Success 201 {object} interface{} "User created successfully"
// @Failure 400 {object} errorHandler.HttpError "Invalid request payload"
// @Failure 422 {object} errorHandler.HttpError "Validation errors"
// @Failure 500 {object} errorHandler.HttpError "Internal server error"
// @Router /api/v1/users [post]
func (h *UserAdmin) Create(ctx *gin.Context) {
	payload := &admin.RequestCreateAdmin{}
	if err := ctx.ShouldBind(payload); err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}

	err := h.adminService.CreateUser(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{})
}
