package producthandler

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/internal/service/productservice"
	"auth-service-rizkysr90-pos/pkg/errorHandler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService productservice.Service
}

func NewCategoryHandler(
	productService productservice.Service,
) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param product body payload.ReqCreateProduct true "Product details"
// @Success 201 {object} payload.ResCreateProduct
// @Failure 400 {object} errorHandler.HttpError
// @Failure 500 {object} errorHandler.HttpError
// @Router /products [post]
func (c *ProductHandler) CreateProduct(ctx *gin.Context) {
	payload := &payload.ReqCreateProduct{}
	if err := ctx.ShouldBind(payload); err != nil {
		err := errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage(err.Error()),
		)
		ctx.Error(err)
		return
	}
	data, err := c.productService.CreateProduct(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, data)
}
// GetAllProducts godoc
// @Summary Get all products
// @Description Retrieve a list of products with optional filtering and pagination
// @Tags products
// @Accept json
// @Produce json
// @Param page_size query int false "Number of items per page" default(10)
// @Param page_number query int false "Page number" default(1)
// @Param category_id query string false "Category ID to filter products"
// @Success 200 {object} payload.ResGetAllProducts
// @Failure 400 {object} errorHandler.HttpError
// @Failure 500 {object} errorHandler.HttpError
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(ctx *gin.Context) {
	request := &payload.ReqGetAllProducts{}
	if err := ctx.ShouldBindQuery(request); err != nil {
		ctx.Error(errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid query parameters"),
			errorHandler.WithMessage(err.Error()),
		))
		return
	}

	response, err := h.productService.GetAllProducts(ctx, request)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product's details
// @Tags products
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param product body payload.ReqUpdateProduct true "Product information to update"
// @Success 200 {object} payload.ResUpdateProduct
// @Failure 400 {object} errorHandler.HttpError
// @Failure 404 {object} errorHandler.HttpError "Product not found"
// @Failure 500 {object} errorHandler.HttpError
// @Router /api/v1/products/{product_id} [put]
func (c *ProductHandler) UpdateProduct(ctx *gin.Context) {
	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.Error(errorHandler.NewBadRequest(
			errorHandler.WithInfo("missing product ID"),
			errorHandler.WithMessage("Product ID is required"),
		))
		return
	}

	payload := &payload.ReqUpdateProduct{}
	if err := ctx.ShouldBindJSON(payload); err != nil {
		ctx.Error(errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid payload"),
			errorHandler.WithMessage(err.Error()),
		))
		return
	}
	// Set the product_id from the URL parameter
	payload.ProductID = productID

	data, err := c.productService.UpdateProduct(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Retrieve a single product's details by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} payload.ResGetProductByID
// @Failure 400 {object} errorHandler.HttpError
// @Failure 404 {object} errorHandler.HttpError "Product not found"
// @Failure 500 {object} errorHandler.HttpError
// @Router /api/v1/products/{product_id} [get]
func (c *ProductHandler) GetProductByID(ctx *gin.Context) {
	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.Error(errorHandler.NewBadRequest(
			errorHandler.WithInfo("missing product ID"),
			errorHandler.WithMessage("Product ID is required"),
		))
		return
	}
	payload := &payload.ReqGetProductByID{}
	// Set the product_id from the URL parameter
	payload.ProductID = productID

	data, err := c.productService.GetProductByID(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// DeleteProductByID godoc
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} payload.ResDeleteProductByID
// @Failure 400 {object} errorHandler.HttpError
// @Failure 404 {object} errorHandler.HttpError
// @Failure 500 {object} errorHandler.HttpError
// @Router /api/v1/products/{product_id} [delete]
func (c *ProductHandler) DeleteProductByID(ctx *gin.Context) {
	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.Error(errorHandler.NewBadRequest(
			errorHandler.WithInfo("missing product ID"),
			errorHandler.WithMessage("Product ID is required"),
		))
		return
	}
	payload := &payload.ReqDeleteProductByID{}
	// Set the product_id from the URL parameter
	payload.ProductID = productID

	data, err := c.productService.DeleteProductByID(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}
