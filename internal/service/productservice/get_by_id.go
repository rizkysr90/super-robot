package productservice

import (
	"context"
	"database/sql"
	"errors"
	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/pkg/errorHandler"
	"strings"

	"github.com/go-playground/validator/v10"
)

type reqGetByID struct {
	data *payload.ReqGetProductByID
}

func (request *reqGetByID) validate() error {
	validate := validator.New()

	// Register custom validation for ProductID if needed
	// For example, if ProductID should be a specific format:
	// validate.RegisterValidation("product_id_format", validateProductIDFormat)

	err := validate.Struct(request.data)
	if err != nil {
		//nolint:errorlint
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return errorHandler.NewInternalServer()
		}
		return errorHandler.NewBadRequest(
			errorHandler.WithInfo("validation error"),
			errorHandler.WithMessage(formatValidationErrors(validationErrors)),
		)
	}

	return nil
}
func formatValidationErrors(errors validator.ValidationErrors) string {
	var errorMessages []string
	for _, err := range errors {
		switch err.Tag() {
		case "required":
			errorMessages = append(errorMessages, err.Field()+" is required")
		// Add more cases here for other validation rules
		default:
			errorMessages = append(errorMessages, err.Field()+" is invalid")
		}
	}
	return "Validation failed: " + strings.Join(errorMessages, ", ")
}
func (request *reqGetByID) sanitize() {
	request.data.ProductID = strings.TrimSpace(request.data.ProductID)
}
func (s *Service) GetProductByID(
	ctx context.Context, request *payload.ReqGetProductByID) (*payload.ResGetProductByID, error) {
	input := reqGetByID{request}
	input.sanitize()
	if err := input.validate(); err != nil {
		return nil, err
	}
	productData, err := s.productStore.GetByID(ctx, input.data.ProductID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorHandler.NewBadRequest(errorHandler.WithInfo("product id not found"))
		}
		return nil, err
	}
	productResponse := payload.ProductData{
		ProductID:     productData.ProductID,
		ProductName:   productData.ProductName,
		Price:         productData.Price,
		BasePrice:     productData.BasePrice,
		StockQuantity: productData.StockQuantity,
		CategoryID:    productData.CategoryID,
		CategoryName:  productData.Category.CategoryName,
		CreatedAt:     productData.CreatedAt,
		UpdatedAt:     productData.UpdatedAt,
		DeletedAt:     productData.DeletedAt,
	}
	return &payload.ResGetProductByID{Data: productResponse}, nil
}
