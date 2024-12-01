package category

import (
	"net/http"
	"strconv"

	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/service"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(
	categoryService service.CategoryService,
) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category with the provided name
// @Tags categories
// @Accept json
// @Produce json
// @Param category body payload.ReqCreateCategory true "Category to create"
// @Success 201 {object} payload.ResCreateCategory
// @Failure 400 {object} errorHandler.HttpError
// @Failure 500 {object} errorHandler.HttpError
// @Router /categories [post]
func (c *CategoryHandler) CreateCategory(ctx *gin.Context) {
	payload := &payload.ReqCreateCategory{}
	if err := ctx.ShouldBind(payload); err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}
	data, err := c.categoryService.Create(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, data)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Retrieve all categories with optional pagination
// @Tags categories
// @Accept json
// @Produce json
// @Param page_number query int false "Page number (default: 1)" minimum(1) default(1)
// @Param page_size query int false "Page size (default: 20)" minimum(1) maximum(100) default(20)
// @Success 200 {object} payload.ResGetAllCategory
// @Failure 400 {object} errorHandler.HttpError
// @Failure 500 {object} errorHandler.HttpError
// @Router /categories [get]
func (c *CategoryHandler) GetAllCategories(ctx *gin.Context) {
	payload := &payload.ReqGetAllCategory{}

	// Extract page_number and page_size from query parameters
	pageNumber := ctx.DefaultQuery("page_number", "1")
	pageSize := ctx.DefaultQuery("page_size", "20")

	// Convert string to int
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid page_number"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid page_size"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}

	// Assign to payload
	payload.PageNumber = pageNumberInt
	payload.PageSize = pageSizeInt
	data, err := c.categoryService.GetCategories(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// GetCategoryByID godoc
// @Summary Get a category by ID
// @Description Retrieve a specific category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param category_id path string true "Category ID"
// @Success 200 {object} payload.ResGetCategoryByID
// @Failure 400 {object} errorHandler.HttpError
// @Failure 404 {object} errorHandler.HttpError
// @Failure 500 {object} errorHandler.HttpError
// @Router /categories/{category_id} [get]
func (c *CategoryHandler) GetCategoryByID(ctx *gin.Context) {
	categoryID := ctx.Param("category_id")
	payloadRequest := payload.ReqGetCategoryByID{
		CategoryID: categoryID,
	}

	data, err := c.categoryService.GetCategoryByID(ctx, &payloadRequest)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// EditCategoryByID godoc
// @Summary Edit a category by ID
// @Description Update a specific category's name by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param category_id path string true "Category ID"
// @Param category body payload.ReqUpdateCategory true "Updated category information"
// @Success 200 {object} payload.ResUpdateCategory
// @Failure 400 {object} errorHandler.HttpError
// @Failure 404 {object} errorHandler.HttpError
// @Failure 500 {object} errorHandler.HttpError
// @Router /categories/{category_id} [put]
func (c *CategoryHandler) EditCategoryByID(ctx *gin.Context) {
	categoryID := ctx.Param("category_id")
	payloadRequest := payload.ReqUpdateCategory{
		ID: categoryID,
	}
	if err := ctx.ShouldBind(&payloadRequest); err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}
	data, err := c.categoryService.EditCategory(ctx, &payloadRequest)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a specific category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param category_id path string true "Category ID"
// @Success 200 {object} payload.ResDeleteCategory
// @Failure 400 {object} errorHandler.HttpError
// @Failure 404 {object} errorHandler.HttpError
// @Failure 500 {object} errorHandler.HttpError
// @Router /categories/{category_id} [delete]
func (c *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	categoryID := ctx.Param("category_id")
	payloadRequest := payload.ReqDeleteCategory{
		ID: categoryID,
	}

	data, err := c.categoryService.DeleteCategory(ctx, &payloadRequest)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}
