package service

import (
	"context"

	payload "auth-service-rizkysr90-pos/internal/payload/http/auth"
)

type AuthService interface {
	CreateUser(ctx context.Context, req *payload.ReqCreateAccount) error
	LoginUser(ctx context.Context,
		req *payload.ReqLoginUser) (*payload.ResLoginUser, error)
	RefreshToken(ctx context.Context,
		req *payload.ReqRefreshToken) (*payload.ResRefreshToken, error)
}
