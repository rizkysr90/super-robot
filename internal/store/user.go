package store

import (
	"context"
	"database/sql"
	"time"
)

type UserData struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}
type InsertedData struct {
	CreatedAt time.Time
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
}
type UserFilterBy struct {
	Email string
}
type UserStore interface {
	Create(ctx context.Context, data *InsertedData) error
	FindOne(ctx context.Context, filterBy *UserFilterBy, staging string) (*UserData, error)
}
