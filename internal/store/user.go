package store

import (
	"context"
	"database/sql"
	"time"
)

type UserData struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}
type InsertedData struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}
type UserFilterBy struct {
	Email string
}
type UserStore interface {
	Create(ctx context.Context, data *InsertedData) error
	FindOne(ctx context.Context, filterBy *UserFilterBy, staging string) (*UserData, error)
}
