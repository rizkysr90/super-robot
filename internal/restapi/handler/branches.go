package handler

import (
	"net/http"
	"rizkysr90-pos/internal/service/branches"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/gin-gonic/gin"
)

type Branches struct {
	branchService *branches.Service
}

func NewBranchHandler(branchService *branches.Service) *Branches {
	return &Branches{branchService: branchService}
}

// @Summary Create Branch
// @Description Creates a new branch with the provided details
// @Tags Branches
// @Accept json
// @Produce json
// @Param payload body branches.RequestCreate true "Branch creation details"
// @Success 201 {object} interface{} "Branch created successfully"
// @Failure 400 {object} errorHandler.HttpError "Invalid request payload"
// @Failure 422 {object} errorHandler.HttpError "Validation errors"
// @Failure 500 {object} errorHandler.HttpError "Internal server error"
// @Router /api/v1/branches [post]
func (h *Branches) Create(ctx *gin.Context) {
	payload := &branches.RequestCreate{}
	if err := ctx.ShouldBind(payload); err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}
	err := h.branchService.Create(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{})
}
