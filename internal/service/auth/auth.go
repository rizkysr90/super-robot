package auth

import (
	"api-iad-ams/internal/service"
	"api-iad-ams/internal/store"
	"context"
)

type AuthService struct {
	userStore store.AuthStore
}

func NewAuthService(userStore store.AuthStore) *AuthService {
	return &AuthService{
		userStore: userStore,
	}
}

func (a *AuthService) CreateUser(ctx context.Context, req *service.CreateUserRequest) error {
	return nil
}
