package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type TenantRole struct {
	db *sql.DB
}

func NewTenantRole(db *sql.DB) *TenantRole {
	return &TenantRole{
		db: db,
	}
}

func (t *TenantRole) Insert(ctx context.Context, data *store.TenantRoleData) error {
	query := `
		INSERT INTO tenant_roles (
			id, 
			tenant_id, 
			name, 
			is_head_office_role, 
			description,
			created_at,
			created_by
		)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	createFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			data.ID,
			data.TenantID,
			data.Name,
			data.IsHeadOfficeRole,
			data.Description,
			data.CreatedAt,
			data.CreatedBy,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}
