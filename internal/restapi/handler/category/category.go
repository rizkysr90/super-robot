package category

import (
	"auth-service-rizkysr90-pos/internal/config"
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/internal/service"
	"auth-service-rizkysr90-pos/pkg/errorHandler"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
	config          config.Config
}

func NewCategoryHandler(
	categoryService service.CategoryService,
) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

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
