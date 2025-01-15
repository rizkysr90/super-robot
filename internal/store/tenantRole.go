package store

import (
	"context"
	"database/sql"
	"time"
)

type TenantRoleData struct {
	ID               string
	TenantID         string
	Name             string
	IsHeadOfficeRole bool
	Description      string
	CreatedAt        time.Time
	UpdatedAt        sql.NullTime
	DeletedAt        sql.NullTime
	CreatedBy        string
	UpdatedBy        sql.NullString
	DeletedBy        sql.NullString
}

type TenantRole interface {
	Insert(ctx context.Context, data *TenantRoleData) error
}
