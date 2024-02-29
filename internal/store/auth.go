package store

import (
	"context"
	"time"
)

type AuthData struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type CreateAccountStore struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}
type AuthStore interface {
	CreateAccount(ctx context.Context, payload *CreateAccountStore) error
}
