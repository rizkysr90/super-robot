package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) Insert(ctx context.Context, userData *store.UserData) error {
	query := `
		INSERT INTO users (
            id,
            email,
            full_name,
            google_id,
            password_hash,
            auth_type,
            user_type,
            tenant_id,
            created_at,
            last_login_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        )
	`
	createFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			userData.ID,
			userData.Email,
			userData.FullName,
			userData.GoogleID,
			userData.PasswordHash,
			userData.AuthType,
			userData.UserType,
			userData.TenantID,
			userData.CreatedAt,
			userData.LastLoginAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}
