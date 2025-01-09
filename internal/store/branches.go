package store

import (
	"context"
	"database/sql"
	"time"
)

type BranchesData struct {
	BranchID       string
	TenantID       string
	Name           string
	Address        string
	CreatedAt      time.Time
	UpdatedAt      sql.NullTime
	DeletedAt      sql.NullTime
	CreatedBy      sql.NullString
	CreatedByOwner sql.NullString
}
type BranchesFilter struct {
	ID       string
	Name     string
	TenantID string
}
type Branches interface {
	FindOne(ctx context.Context, filter *BranchesFilter) (*BranchesData, error)
	Insert(ctx context.Context, branchData *BranchesData) error
	TotalBranches(ctx context.Context, tenantID string) (uint8, error)
}
