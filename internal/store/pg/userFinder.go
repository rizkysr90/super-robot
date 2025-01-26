package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/constant"
	"rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type UserFinder struct {
	db *sql.DB
}

func NewUserFinder(db *sql.DB) *UserFinder {
	return &UserFinder{
		db: db,
	}
}

func (u *UserFinder) FindOne(ctx context.Context,
	filter *store.UserQueryFilter) (*store.UserData, error) {
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
