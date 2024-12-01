package service

import (
	"context"
	"rizkysr90-pos/internal/payload"
)

type CategoryService interface {
	Create(ctx context.Context,
		req *payload.ReqCreateCategory) (*payload.ResCreateCategory, error)
	GetCategories(ctx context.Context,
		request *payload.ReqGetAllCategory) (*payload.ResGetAllCategory, error)
	GetCategoryByID(ctx context.Context,
		request *payload.ReqGetCategoryByID) (*payload.ResGetCategoryByID, error)
	EditCategory(ctx context.Context, request *payload.ReqUpdateCategory) (*payload.ResUpdateCategory, error)
	DeleteCategory(ctx context.Context,
		request *payload.ReqDeleteCategory) (*payload.ResDeleteCategory, error)
}
