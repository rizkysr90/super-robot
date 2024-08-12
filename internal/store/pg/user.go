package pg

import (
	"context"
	"database/sql"
	"time"

	"auth-service-rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type Users struct {
	db *sql.DB
}

func NewUserDB(db *sql.DB) *Users {
	return &Users{
		db: db,
	}
}

func(u *Users) FindByUsername(ctx context.Context, username string) (*store.User, error) {
	query := `
		SELECT id, username, password, access_token FROM users
		WHERE username = $1 AND deleted_at IS NULL
	`
	data := &store.User{}
	row := sqldb.WithinTxContextOrDB(ctx, u.db).
		QueryRowContext(ctx, query, username)
	if err := row.Err(); err != nil {
		return nil, err
	}
	err := row.Scan(&data.ID, &data.Username,&data.Password, &data.AccessToken)
	if err != nil {
		return nil, err
	}
	return data, nil

}

func (u *Users) Create(ctx context.Context, user *store.User) error {
	query := `
		INSERT INTO users (id, username, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	createFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			user.ID,
			user.Username,
			user.Password,
			user.CreatedAt,
			user.UpdatedAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}

func (u *Users) UpdateUserAccessToken(ctx context.Context, username, token string) error {
	query := `
		UPDATE users SET access_token = $1, updated_at = $2 WHERE username = $3
	`
	updateFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			token,
			time.Now().UTC(),
			username,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, updateFunc)
}