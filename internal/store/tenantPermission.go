package store

import (
	"context"
	"database/sql"
	"time"
)

type TenantPermissionData struct {
	ID             string
	PermissionCode string
	TenantRoleID   string
	TenantID       string
	CreatedAt      time.Time
	UpdatedAt      sql.NullTime
	DeletedAt      sql.NullTime
	CreatedBy      string
	UpdatedBy      sql.NullString
	DeletedBy      sql.NullString
}

type TenantPermission interface {
	FindByTenantRoleID(ctx context.Context, tenantRoleIds []string) ([]TenantPermissionData, error)
}
