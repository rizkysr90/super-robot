package store

import (
	"context"
	"database/sql"
	"time"
)

type WorkLocationData struct {
	ID        string
	TenantID  string
	UserID    string
	BranchID  sql.NullString
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
	CreatedBy string
	UpdatedBy sql.NullString
	DeletedBy sql.NullString
}

type WorkLocation interface {
	FindByUserID(ctx context.Context, userID string) ([]WorkLocationData, error)
}
