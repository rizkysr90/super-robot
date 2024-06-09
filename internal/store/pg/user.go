package pg

import (
	"context"
	"database/sql"

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

func (u *Users) Create(ctx context.Context, data *store.User) error {
	createFunc := func(tx sqldb.QueryExecutor) error {
		query := `
			INSERT INTO users 
			(id, first_name, last_name,
			 password, email, phone, 
			 is_activated, role, created_at
			)
			VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
		_, err := tx.ExecContext(ctx, query,
			data.ID,
			data.FirstName,
			data.LastName,
			data.Password,
			data.Email,
			data.Phone,
			data.IsActivated,
			data.Role,
			data.CreatedAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}

const FindActiveUserByEmailSQL = `
	SELECT id 
	FROM users
	WHERE 
	is_activated = true AND
	deleted_at IS NULL AND
	email = $1
`

func (u *Users) FindActiveUserByEmail(ctx context.Context, filter *store.User) (*store.User, error) {
	resultData := &store.User{}
	row := sqldb.WithinTxContextOrDB(ctx, u.db).
		QueryRowContext(ctx, FindActiveUserByEmailSQL, filter.Email)
	if err := row.Err(); err != nil {
		return nil, err
	}
	row.Scan(&resultData.ID)

	return resultData, nil
}
