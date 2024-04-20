package service

import (
	payload "auth-service-rizkysr90-pos/internal/payload/http/store"
	"context"
)

type StoreService interface {
	CreateStore(ctx context.Context, req *payload.ReqCreateStore) (*payload.ResCreateStore, error)
	GetAllStore(ctx context.Context, req *payload.ReqGetAllStore) (*payload.ResGetAllStore, error)
	DeleteStore(ctx context.Context, req *payload.ReqDeleteStore) error
	UpdateStore(ctx context.Context, req *payload.ReqUpdateStore) error
}
