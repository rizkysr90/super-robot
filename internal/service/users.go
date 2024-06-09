package service

import (
	"context"

	payload "auth-service-rizkysr90-pos/internal/payload/http/users"
)

type UsersService interface {
	CreateUser(ctx context.Context, req *payload.ReqCreateUsers) error
	// LoginUser(ctx context.Context,
	// 	req *payload.ReqLoginUser) (*payload.ResLoginUser, error)
	// RefreshToken(ctx context.Context,
	// 	req *payload.ReqRefreshToken) (*payload.ResRefreshToken, error)
}
