package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
	ID          string
	Username    string
	Password    string
	AccessToken sql.NullString
}
type UserStore interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	UpdateUserAccessToken(ctx context.Context, username, token string) error
}
