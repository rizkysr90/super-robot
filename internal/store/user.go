package store

import (
	"context"
	"database/sql"
	"time"
)

type UserData struct {
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
	DeletedAt   sql.NullTime
	ID          string
	FirstName   string
	LastName    string
	Email       string
	Password    string
	IsActivated bool
}
type InsertedData struct {
	CreatedAt   time.Time
	ID          string
	FirstName   string
	LastName    string
	Email       string
	Password    string
	IsActivated bool
}
type UserFilterBy struct {
	Email string
}
type UserStore interface {
	Create(ctx context.Context, data *InsertedData) error
	FindOne(ctx context.Context, filterBy *UserFilterBy, staging string) (*UserData, error)
}
