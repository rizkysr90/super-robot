package service

import (
	"context"

	payload "github.com/rizkysr90/go-boilerplate/internal/payload/http/auth"
)

type AuthService interface {
	CreateUser(ctx context.Context, req *payload.ReqCreateAccount) error
}
