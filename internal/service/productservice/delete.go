package productservice

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/pkg/errorHandler"
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type reqDeleteByID struct {
	data *payload.ReqDeleteProductByID
}

func (request *reqDeleteByID) validate() error {
	validate := validator.New()

	err := validate.Struct(request.data)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return errorHandler.NewBadRequest(
			errorHandler.WithInfo("validation error"),
			errorHandler.WithMessage(formatValidationErrors(validationErrors)),
		)
	}

	return nil
}

func (request *reqDeleteByID) sanitize() {
	request.data.ProductID = strings.TrimSpace(request.data.ProductID)
}

func (s *Service) DeleteProductByID(ctx context.Context, request *payload.ReqDeleteProductByID) (*payload.ResDeleteProductByID, error) {
	input := reqDeleteByID{data: request}
	input.sanitize()
	if err := input.validate(); err != nil {
		return nil, err
	}

	err := sqldb.WithinTx(ctx, s.db, func(qe sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, qe)
		return s.productStore.DeleteByID(txContext, input.data.ProductID)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorHandler.NewNotFound(errorHandler.WithInfo("product not found"))
		}
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo("failed to delete product"),
			errorHandler.WithMessage(err.Error()),
		)
	}

	return &payload.ResDeleteProductByID{
	}, nil
}