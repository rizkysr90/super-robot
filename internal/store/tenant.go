package store

import (
	"context"
	"database/sql"
	"time"
)

type TenantData struct {
	CreatedAt time.Time
	Owner     *UserData
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
	ID        string
	Name      string
	OwnerID   sql.NullString
}

type Tenant interface {
	Insert(ctx context.Context, tenantData *TenantData) error
	Update(ctx context.Context, tenantData *TenantData) error
}
