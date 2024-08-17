package service

import (
	"context"

	"auth-service-rizkysr90-pos/internal/payload"
)

type UsersService interface {
	LoginUser(ctx context.Context, req *payload.ReqLoginUsers) (*payload.ResLoginUsers, error)
	CreateUser(ctx context.Context, req *payload.ReqCreateUsers) (*payload.ResCreateUsers, error)
	// 	req *payload.ReqLoginUser) (*payload.ResLoginUser, error)
	// RefreshToken(ctx context.Context,
	// 	req *payload.ReqRefreshToken) (*payload.ResRefreshToken, error)
}
