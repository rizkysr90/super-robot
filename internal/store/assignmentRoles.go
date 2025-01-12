package store

import (
	"context"
	"database/sql"
	"time"
)

type AssignmentRoleData struct {
	ID           string
	TenantID     string
	AssignmentID string
	TenantRoleID string
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
	DeletedAt    sql.NullTime
	CreatedBy    string
	UpdatedBy    sql.NullString
	DeletedBy    sql.NullString
}

type AssignmentRole interface {
	FindMany(ctx context.Context,
		assignmentIDs []string) ([]AssignmentRoleData, error)
}
