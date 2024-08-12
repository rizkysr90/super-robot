package category

import (
	"auth-service-rizkysr90-pos/internal/config"
	payload "auth-service-rizkysr90-pos/internal/payload/http/category"
	"auth-service-rizkysr90-pos/internal/service"
	"auth-service-rizkysr90-pos/pkg/errorHandler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
	config       config.Config
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
	data, err := c.categoryService.Create(ctx, payload);
	if  err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}