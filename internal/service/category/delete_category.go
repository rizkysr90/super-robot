package category

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/pkg/errorHandler"
	"auth-service-rizkysr90-pos/pkg/validator"
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type requestDeleteCategory struct {
	*payload.ReqDeleteCategory
}

func (request *requestDeleteCategory) sanitize() {
	request.ID = strings.TrimSpace(request.ID)
}
func (request *requestDeleteCategory) validate() error {
	if !validator.ValidateRequired(request.ID, "id") {
		return errorHandler.NewBadRequest(errorHandler.WithInfo("required id data"))
	}
	return nil
}
func (c *Service) DeleteCategory(ctx context.Context,
	request *payload.ReqDeleteCategory) (*payload.ResDeleteCategory, error) {
	input := &requestDeleteCategory{request}

	input.sanitize()
	if err := input.validate(); err != nil {
		return nil, err
	}
	_, err := c.categoryStore.FindByID(ctx, input.ID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorHandler.NewNotFound(errorHandler.WithInfo("category not found"))
		}
		return nil, err
	}
	err = sqldb.WithinTx(ctx, c.db, func(qe sqldb.QueryExecutor) error {
		txCtx := sqldb.WithTxContext(ctx, qe)
		return c.categoryStore.SoftDelete(txCtx, input.ID)
	})
	if err != nil {
		return nil, err
	}
	return &payload.ResDeleteCategory{}, nil
}
