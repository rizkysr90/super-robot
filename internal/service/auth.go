package service

import (
	payload "api-iad-ams/internal/payload/http/auth"
	"context"
)

type AuthService interface {
	CreateUser(ctx context.Context, req *payload.ReqCreateAccount) error
}
