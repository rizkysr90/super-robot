package service

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"context"
)


type CategoryService interface {
	Create(ctx context.Context,
		req *payload.ReqCreateCategory) (*payload.ResCreateCategory, error)
}