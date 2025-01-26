package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type TenantFinder struct {
	db *sql.DB
}

func NewTenantFinder(db *sql.DB) *TenantFinder {
	return &TenantFinder{
		db: db,
	}
}

func (t *TenantFinder) FindOne(ctx context.Context,
	filter *store.TenantFilter) (*store.TenantData, error) {
	query := `
		SELECT id, name, owner_id, created_at FROM tenants
		WHERE 
			$1 = '' OR id = $1::uuid AND
			$2 = '' OR name = $2 AND
		deleted_at IS NULL
	`
	data := &store.TenantData{}
	row := sqldb.WithinTxContextOrDB(ctx, t.db).
		QueryRowContext(ctx, query, filter.ID, filter.Name)
	if err := row.Err(); err != nil {
		return nil, err
	}
	err := row.Scan(&data.ID, &data.Name,
		&data.OwnerID, &data.CreatedAt)
	if err != nil {
		return nil, err
	}
	return data, nil
}
