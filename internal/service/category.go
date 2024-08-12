package service

import (
	payload "auth-service-rizkysr90-pos/internal/payload/http/category"
	"context"
)


type CategoryService interface {
	Create(ctx context.Context,
		req *payload.ReqCreateCategory) (*payload.ResCreateCategory, error)
}