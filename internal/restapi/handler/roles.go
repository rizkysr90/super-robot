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
