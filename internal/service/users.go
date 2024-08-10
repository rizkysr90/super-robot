package service

import (
	"context"

	 "auth-service-rizkysr90-pos/internal/payload/http/users"
)

type UsersService interface {
	LoginUser(ctx context.Context, user)
	// 	req *payload.ReqLoginUser) (*payload.ResLoginUser, error)
	// RefreshToken(ctx context.Context,
	// 	req *payload.ReqRefreshToken) (*payload.ResRefreshToken, error)
}
