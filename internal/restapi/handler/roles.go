package handler

import (
	"net/http"
	"rizkysr90-pos/internal/service/roles"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/gin-gonic/gin"
)

type Roles struct {
	rolesService *roles.Service
}

func NewRolesHandler(rolesService *roles.Service) *Roles {
	return &Roles{rolesService: rolesService}
}

// @Summary Create Role
// @Description Creates a new role with the provided details
// @Tags Roles
// @Accept json
// @Produce json
// @Param payload body roles.RequestCreateRoles true "Role creation request"
// @Success 201 {object} interface{} "Role created successfully"
// @Failure 400 {object} errorHandler.HttpError "Invalid request payload"
// @Failure 422 {object} errorHandler.HttpError "Validation errors"
// @Failure 500 {object} errorHandler.HttpError "Internal server error"
// @Router /api/v1/roles [post]
func (h *Roles) Create(ctx *gin.Context) {
	payload := &roles.RequestCreateRoles{}
	if err := ctx.ShouldBind(payload); err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}
	err := h.rolesService.Create(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{})
}
