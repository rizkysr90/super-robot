package productservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"rizkysr90-pos/internal/helper"
	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/go-playground/validator/v10"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

// Update Product

type reqUpdateProduct struct {
	data *payload.ReqUpdateProduct
}

func (req *reqUpdateProduct) sanitize() {
	req.data.ProductName = strings.TrimSpace(req.data.ProductName)
	req.data.CategoryID = strings.TrimSpace(req.data.CategoryID)
	req.data.ProductName = strings.ToUpper(req.data.ProductName)
}

func (req *reqUpdateProduct) validate() error {
	httpErrors := []errorHandler.HttpError{}
	validate := validator.New()
	err := validate.Struct(req.data)
	if err != nil {
		//nolint:errorlint
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return errorHandler.NewInternalServer()
		}
		for _, e := range validationErrors {
			httpError := errorHandler.HttpError{
				Code:    400,
				Info:    "Validation Error",
				Message: helper.FormatErrorMessage(e),
			}
			httpErrors = append(httpErrors, httpError)
		}
	}
	if len(httpErrors) > 0 {
		return errorHandler.NewMultipleFieldsValidation(httpErrors)
	}
	return nil
}

func (s *Service) UpdateProduct(
	ctx context.Context, request *payload.ReqUpdateProduct) (*payload.ResUpdateProduct, error) {
	input := &reqUpdateProduct{request}
	input.sanitize()
	if err := input.validate(); err != nil {
		return nil, err
	}
	product, err := s.productStore.GetByName(ctx, input.data.ProductName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if product != nil && product.ProductID != input.data.ProductID {
		errMsg := fmt.Sprintf("duplicate product name : %s", input.data.ProductName)
		return nil, errorHandler.NewBadRequest(errorHandler.WithInfo(errMsg))
	}
	updateData := store.ProductData{
		ProductID:     request.ProductID,
		ProductName:   request.ProductName,
		Price:         request.Price,
		BasePrice:     request.BasePrice,
		StockQuantity: request.StockQuantity,
		CategoryID:    request.CategoryID,
	}

	if err = sqldb.WithinTx(ctx, s.db, func(qe sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, qe)
		return s.productStore.Update(txContext, &updateData)
	}); err != nil {
		return nil, err
	}

	return &payload.ResUpdateProduct{}, nil
}
