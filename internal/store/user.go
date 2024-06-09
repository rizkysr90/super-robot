package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
	DeletedAt    sql.NullTime
	ID           string
	FirstName    string
	LastName     string
	Phone        string
	Email        string
	Password     string
	RefreshToken sql.NullString
	Role         int
	IsActivated  bool
}
type UserStore interface {
	Create(ctx context.Context, createdUser *User) error
	FindActiveUserByEmail(ctx context.Context, filter *User) (*User, error)
}
