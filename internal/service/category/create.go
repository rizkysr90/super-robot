package category

import (
	"context"
	"database/sql"
	"errors"
	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store"
	"strings"
	"time"

	"rizkysr90-pos/pkg/errorHandler"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type reqCreateCategory struct {
	*payload.ReqCreateCategory
}

func (req *reqCreateCategory) sanitize() {
	req.CategoryName = strings.TrimSpace(req.CategoryName)
	req.CategoryName = strings.ToUpper(req.CategoryName)
}
func (req *reqCreateCategory) validate() error {
	if len(req.CategoryName) == 0 {
		return errorHandler.NewBadRequest(errorHandler.WithInfo("category name is required"))
	}
	if len(req.CategoryName) > 100 {
		return errorHandler.NewBadRequest(errorHandler.WithInfo("max category name is 100 characters"))
	}
	return nil
}
func (c *Service) Create(ctx context.Context,
	req *payload.ReqCreateCategory) (*payload.ResCreateCategory, error) {
	input := &reqCreateCategory{req}
	input.sanitize()
	if err := input.validate(); err != nil {
		return nil, err
	}
	category, err := c.categoryStore.FindByName(ctx, req.CategoryName)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return nil, err
	}
	if category != nil {
		return nil, errorHandler.NewBadRequest(errorHandler.WithInfo("duplicate category name"))
	}
	insertedData := &store.CategoryData{
		ID:           uuid.NewString(),
		CategoryName: input.CategoryName,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	err = sqldb.WithinTx(ctx, c.db, func(qe sqldb.QueryExecutor) error {
		tx := sqldb.WithTxContext(ctx, qe)
		return c.categoryStore.Create(tx, insertedData)
	})
	if err != nil {
		return nil, err
	}
	return &payload.ResCreateCategory{}, nil
}
