package pg

import (
	"api-iad-ams/internal/store"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth struct {
	db *pgxpool.Pool
}

func NewAuthDB(db *pgxpool.Pool) *Auth {
	return &Auth{
		db: db,
	}
}

func (u *Auth) CreateAccount(ctx context.Context, payload *store.CreateAccountStore) error {
	return nil
}
