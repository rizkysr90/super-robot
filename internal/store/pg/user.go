package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rizkysr90/go-boilerplate/internal/store"
	"github.com/rizkysr90/go-boilerplate/pkg/sqldb"
)

type User struct {
	db *sql.DB
}

func NewUserDB(db *sql.DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) Create(ctx context.Context, data *store.InsertedData) error {
	createFunc := func(tx sqldb.QueryExecutor) error {
		query := `
			INSERT INTO users 
			(id, first_name, last_name,
			 password, email, created_at, is_activated
			)
			VALUES 
			($1, $2, $3, $4, $5, $6, $7)
		`
		_, err := tx.ExecContext(ctx, query,
			data.ID,
			data.FirstName,
			data.LastName,
			data.Password,
			data.Email,
			data.CreatedAt,
			data.IsActivated,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}
func (u *User) FindOne(ctx context.Context,
	filterBy *store.UserFilterBy, staging string) (
	*store.UserData, error) {
	var result store.UserData
	switch staging {
	case "findactiveuser":
		// 1 is for filter by email
		query := `SELECT email,password
			 FROM users WHERE email = $1 AND is_activated = true`
		err := sqldb.WithinTxContextOrDB(ctx, u.db).
			QueryRowContext(ctx, query, filterBy.Email).
			Scan(&result.Email, &result.Password)
		if err != nil {
			return nil, err
		}
		return &result, nil
	default:
		return nil, errors.New("staging db not found")
	}
}
