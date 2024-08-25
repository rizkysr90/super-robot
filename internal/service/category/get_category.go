package category

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/pkg/errorHandler"
	"auth-service-rizkysr90-pos/pkg/validator"
	"context"
	"database/sql"
	"errors"
	"strings"
)

type reqGetCategoryByID struct {
	*payload.ReqGetCategoryByID
}

func (request *reqGetCategoryByID) sanitize() {
	request.CategoryID = strings.TrimSpace(request.CategoryID)
}
func (request *reqGetCategoryByID) validate() error {
	if !validator.ValidateRequired(request.CategoryID, "category_id") {
		return errorHandler.NewBadRequest(errorHandler.WithInfo("category_id is required"))
	}
	return nil
}
func (c *Service) GetCategoryByID(ctx context.Context,
	request *payload.ReqGetCategoryByID) (*payload.ResGetCategoryByID, error) {
	input := reqGetCategoryByID{request}
	input.sanitize()
	if err := input.validate(); err != nil {
		return nil, err
	}
	category, err := c.categoryStore.FindByID(ctx, request.CategoryID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorHandler.NewNotFound()
		}
	}
	return &payload.ResGetCategoryByID{
		CategoryData: &payload.CategoryData{
			ID:           category.ID,
			CategoryName: category.CategoryName,
			CreatedAt:    category.CreatedAt,
			UpdatedAt:    category.UpdatedAt,
			DeletedAt:    category.DeletedAt.Time,
		},
	}, nil
}
