package pg

import (
	"context"
	"database/sql"
	"errors"

	"auth-service-rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
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
		query := `SELECT id, email,password
			 FROM users WHERE email = $1 AND is_activated = true`
		err := sqldb.WithinTxContextOrDB(ctx, u.db).
			QueryRowContext(ctx, query, filterBy.Email).
			Scan(&result.ID, &result.Email, &result.Password)
		if err != nil {
			return nil, err
		}
		return &result, nil
	default:
		return nil, errors.New("staging db not found")
	}
}

func updateRefreshtoken(ctx context.Context,
	updatedData *store.UserData, f *store.UserFilterBy) func(tx sqldb.QueryExecutor) error {

	result := func(tx sqldb.QueryExecutor) error {
		query := `
			UPDATE users SET refresh_token = $1
			WHERE email = $2
		`
		_, err := tx.ExecContext(ctx, query, updatedData.RefreshToken, f.Email)
		if err != nil {
			return err
		}
		return nil
	}
	return result
}
func (u *User) Update(ctx context.Context,
	updatedData *store.UserData,
	filter *store.UserFilterBy,
	staging string) error {

	var updateFunc func(tx sqldb.QueryExecutor) error
	switch staging {
	case "updaterefreshtoken":
		updateFunc = updateRefreshtoken(ctx, updatedData, filter)
	default:
		return errors.New("staging db not found")
	}
	return sqldb.WithinTxContextOrError(ctx, updateFunc)

}
