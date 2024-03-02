package pg

import (
	"api-iad-ams/internal/store"
	"api-iad-ams/pkg/sqldb"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	db *pgxpool.Pool
}

func NewUserDB(db *pgxpool.Pool) *User {
	return &User{
		db: db,
	}
}

func (u *User) Create(ctx context.Context, payload *store.UserData) error {
	createFunc := func(tx sqldb.QueryExecutor) error {
		query := `
			INSERT INTO users 
			(id, first_name, last_name,
			 password, email, created_at
			)
			VALUES 
			($1, $2, $3, $4, $5, $6)
		`
		_, err := tx.Exec(ctx, query,
			payload.Id,
			payload.FirstName,
			payload.LastName,
			payload.Password,
			payload.Email,
			payload.CreatedAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}
func (u *User) FindOne(ctx context.Context,
	filterBy *store.UserData, staging string) (
	*store.UserData, error) {
	var result store.UserData
	switch staging {
	case "findbyemail":
		// 1 is for filter by email
		query := `SELECT email
			 FROM users WHERE email = $1`
		err := sqldb.WithinTxContextOrDB(ctx, u.db).
			QueryRow(ctx, query, filterBy.Email).
			Scan(&result.Email)
		if err != nil {
			return nil, err
		}
		return &result, nil
	default:
		return nil, errors.New("staging db not found")
	}
}
