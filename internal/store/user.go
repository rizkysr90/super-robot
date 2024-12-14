package store

import (
	"context"
	"database/sql"
	"time"
)

type UserData struct {
	LastLoginAt  time.Time
	CreatedAt    time.Time
	Tenant       *TenantData
	DeletedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	UserType     string
	ID           string
	TenantID     string
	AuthType     string
	FullName     string
	Email        string
	PasswordHash sql.NullString
	GoogleID     sql.NullString
}
type UserQueryFilter struct {
	Email string
}
type User interface {
	Insert(ctx context.Context, userData *UserData) error
	FindOne(ctx context.Context, filter *UserQueryFilter) (*UserData, error)
	Update(ctx context.Context, userData *UserData) error
}
