package category

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/internal/store"
	"auth-service-rizkysr90-pos/internal/utility"
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"auth-service-rizkysr90-pos/pkg/errorHandler"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type reqCreateCategory struct {
	*payload.ReqCreateCategory
}

func (req *reqCreateCategory) sanitize() {
	req.CategoryName = utility.SanitizeReqBody(req.CategoryName)
	req.CategoryName = strings.ToUpper(req.CategoryName)
}
func (req *reqCreateCategory) validate() error {
	if len(req.CategoryName) == 0 {
		return errorHandler.NewBadRequest(errorHandler.WithInfo("category name is required"))
	}
	if len(req.CategoryName) > 100 {
		return errorHandler.NewBadRequest(errorHandler.WithMessage("max category name is 100 characters"))
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
		return nil, errorHandler.NewBadRequest(errorHandler.WithMessage("duplicate category name"))
	}
	insertedData := &store.CategoryData{
		ID:           uuid.NewString(),
		CategoryName: input.CategoryName,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	log.Println("HEREEE", insertedData)
	err = sqldb.WithinTx(ctx, c.db, func(qe sqldb.QueryExecutor) error {
		tx := sqldb.WithTxContext(ctx, qe)
		return c.categoryStore.Create(tx, insertedData)
	})
	if err != nil {
		return nil, err
	}
	return &payload.ResCreateCategory{}, nil
}
