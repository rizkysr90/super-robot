package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type Tenant struct {
	db *sql.DB
}

func NewTenant(db *sql.DB) *Tenant {
	return &Tenant{
		db: db,
	}
}

func (t *Tenant) Insert(ctx context.Context, tenantData *store.TenantData) error {
	query := `
		INSERT INTO tenants (id, name, owner_id, created_at)
        VALUES ($1, $2, $3, $4)
	`
	createFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			tenantData.ID,
			tenantData.Name,
			tenantData.OwnerID,
			tenantData.CreatedAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}
func (t *Tenant) Update(ctx context.Context, tenantData *store.TenantData) error {
	query := `
		UPDATE tenants SET owner_id = $2, updated_at = $3 WHERE id = $1`
	updateFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			tenantData.ID,
			tenantData.OwnerID,
			tenantData.UpdatedAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, updateFunc)
}
