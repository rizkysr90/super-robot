package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/constant"
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
			last_login_at,
			created_by,
			is_verified
		) VALUES (
			$1, 
			$2, 
			$3, 
			NULLIF($4, ''), 
			$5, 
			$6, 
			$7, 
			$8, 
			$9, 
			$10,
			NULLIF($11, ''),
			$12
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
			userData.CreatedBy.String,
			userData.IsVerified,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}

func (u *User) FindOne(ctx context.Context, filter *store.UserQueryFilter) (*store.UserData, error) {
	if filter.ID == "" {
		filter.ID = constant.EmptyUUID
	}
	query := `
		SELECT id, email, full_name, google_id, auth_type, user_type, tenant_id
		FROM users 
		WHERE $1 = '' OR email = $1 AND
		$2 = '00000000-0000-0000-0000-000000000000' OR id = $2::uuid AND 
		deleted_at IS NULL
	`
	data := &store.UserData{}
	row := sqldb.WithinTxContextOrDB(ctx, u.db).
		QueryRowContext(ctx, query, filter.Email, filter.ID)
	if err := row.Err(); err != nil {
		return nil, err
	}
	err := row.Scan(&data.ID, &data.Email, &data.FullName, &data.GoogleID, &data.AuthType,
		&data.UserType, &data.TenantID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *User) Update(ctx context.Context, userData *store.UserData) error {
	query := `
		UPDATE users SET
			email = $2,
			full_name = $3,
			google_id = $4, 
			last_login_at = $5
		WHERE id = $1
	`
	updateFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			userData.ID,
			userData.Email,
			userData.FullName,
			userData.GoogleID,
			userData.LastLoginAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, updateFunc)
}
