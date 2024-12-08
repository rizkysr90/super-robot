package store

import (
	"context"
	"database/sql"
	"time"
)

type TenantData struct {
	ID        string
	Name      string
	OwnerID   sql.NullString
	Owner     *UserData
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

type Tenant interface {
	Insert(ctx context.Context, tenantData *TenantData) error
	Update(ctx context.Context, tenantData *TenantData) error
}
