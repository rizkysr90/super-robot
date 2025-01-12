package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/lib/pq"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type TenantPermission struct {
	db *sql.DB
}

func NewTenantPermission(db *sql.DB) *TenantPermission {
	return &TenantPermission{
		db: db,
	}
}

func (t *TenantPermission) FindByTenantRoleID(ctx context.Context, tenantRoleIds []string) ([]store.TenantPermissionData, error) {
	query := `
		SELECT id, permission_code, tenant_role_id, tenant_id,
			created_at, updated_at, deleted_at,
			created_by, updated_by, deleted_by
		FROM tenant_role_permissions
		WHERE tenant_role_id = ANY($1::UUID[]) AND deleted_at IS NULL
	`
	rows, err := sqldb.WithinTxContextOrDB(ctx, t.db).QueryContext(ctx, query,
		pq.Array(tenantRoleIds),
	)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return nil, err
	}
	roles := make([]store.TenantPermissionData, 0)
	for rows.Next() {
		data := store.TenantPermissionData{}
		err = rows.Scan(
			&data.ID,
			&data.PermissionCode,
			&data.TenantRoleID,
			&data.TenantID,
			&data.CreatedAt,
			&data.UpdatedAt,
			&data.DeletedAt,
			&data.CreatedBy,
			&data.UpdatedBy,
			&data.DeletedBy,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, data)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
