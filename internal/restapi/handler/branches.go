package handler

import (
	"net/http"
	"rizkysr90-pos/internal/service/branches"
	"rizkysr90-pos/pkg/errorHandler"
	"strconv"

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

// @Summary Get Branches
// @Description Retrieves a paginated list of branches with optional filters
// @Tags Branches
// @Accept json
// @Produce json
// @Param page_size query int false "Page size (default: 20)" minimum(1) maximum(100)
// @Param page_number query int false "Page number (default: 1)" minimum(1)
// @Param branch_name query string false "Branch name filter" maxLength(100) pattern(^[a-zA-Z0-9-_]+$)
// @Param tenant_id query string false "Tenant ID filter" maxLength(500)
// @Param action_by query string true "User performing the action" maxLength(500)
// @Success 200 {object} branches.ResponseGetBranches "Successfully retrieved branches"
// @Failure 400 {object} errorHandler.HttpError "Invalid request parameters"
// @Failure 422 {object} errorHandler.HttpError "Validation errors"
// @Failure 500 {object} errorHandler.HttpError "Internal server error"
// @Router /api/v1/branches [get]
func (h *Branches) Branches(ctx *gin.Context) {
	// Initialize the pointer to the struct
	request := &branches.RequestGetBranches{}

	// Parse page_size
	if pageSize := ctx.Query("page_size"); pageSize != "" {
		parsedPageSize, err := strconv.Atoi(pageSize)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorHandler.HttpError{
				Code:    http.StatusBadRequest,
				Message: "Invalid page_size parameter",
			})
			return
		}
		request.PageSize = parsedPageSize
	}

	// Parse page_number
	if pageNumber := ctx.Query("page_number"); pageNumber != "" {
		parsedPageNumber, err := strconv.Atoi(pageNumber)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorHandler.HttpError{
				Code:    http.StatusBadRequest,
				Message: "Invalid page_number parameter",
			})
			return
		}
		request.PageNumber = parsedPageNumber
	}

	// Parse other query parameters
	request.BranchName = ctx.Query("branch_name")
	request.TenantID = ctx.Query("tenant_id")
	request.ActionBy = ctx.Query("action_by")
	response := &branches.ResponseGetBranches{}
	// TODO: Call your service layer
	// response, err := service.GetBranches(ctx.Request.Context(), request)
	// if err != nil {
	//     // Handle error
	//     return
	// }

	ctx.JSON(http.StatusOK, response)
}
