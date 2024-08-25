package category

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/internal/store"
	"auth-service-rizkysr90-pos/pkg/errorHandler"
	"auth-service-rizkysr90-pos/pkg/validator"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type reqEditCategory struct {
	*payload.ReqUpdateCategory
}

func (request *reqEditCategory) sanitize() {
	request.ID = strings.TrimSpace(request.ID)
	request.CategoryName = strings.ToUpper(strings.TrimSpace(request.CategoryName))
}
func (request *reqEditCategory) validate() error {
	if !validator.ValidateRequired(request.ID, "id") {
		return errorHandler.NewBadRequest(errorHandler.WithInfo("id is required"))
	}
	if !validator.ValidateRequired(request.CategoryName, "category name") {
		return errorHandler.NewBadRequest(errorHandler.WithInfo("category name is required"))
	}
	if len(request.CategoryName) > 200 {
		return errorHandler.NewBadRequest(errorHandler.WithInfo("too long, max 200 character"))
	}

	return nil
}

func (c *Service) EditCategory(ctx context.Context,
	request *payload.ReqUpdateCategory) (*payload.ResUpdateCategory, error) {
	input := reqEditCategory{request}
	input.sanitize()
	if err := input.validate(); err != nil {
		return nil, err
	}
	category, err := c.categoryStore.FindByName(ctx, input.CategoryName)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return nil, err
	}
	if category != nil && category.ID != input.ID {
		return nil, errorHandler.NewBadRequest(errorHandler.WithInfo("duplicate category name"))
	}
	updatedData := store.CategoryData{
		ID:           input.ID,
		CategoryName: input.CategoryName,
		UpdatedAt:    time.Now().UTC(),
	}
	if err = sqldb.WithinTx(ctx, c.db, func(qe sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, qe)
		return c.categoryStore.Update(txContext, &updatedData)
	}); err != nil {
		return nil, err
	}
	return &payload.ResUpdateCategory{}, nil
}
