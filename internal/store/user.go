package store

import (
	"context"
	"database/sql"
	"time"
)

type UserData struct {
	ID           string
	Email        string
	FullName     string
	GoogleID     sql.NullString
	PasswordHash sql.NullString
	AuthType     string
	UserType     string
	TenantID     string
	Tenant       *TenantData
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
	DeletedAt    sql.NullTime
	LastLoginAt  time.Time
}
type UserQueryFilter struct {
	Email string
}
type User interface {
	Insert(ctx context.Context, userData *UserData) error
	FindOne(ctx context.Context, filter *UserQueryFilter) (*UserData, error)
	Update(ctx context.Context, userData *UserData) error
}
