package service

import (
	payload "auth-service-rizkysr90-pos/internal/payload/http/employee"
	"context"
)

type EmployeeService interface {
	Create(ctx context.Context, req *payload.ReqCreateEmployee) error
	Login(ctx context.Context, req *payload.ReqLoginEmployee) (*payload.ResLoginEmployee, error)
	RefreshToken(ctx context.Context,
		req *payload.ReqRefreshToken) (*payload.ResRefreshToken, error)
}
