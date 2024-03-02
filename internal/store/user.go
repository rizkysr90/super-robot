package store

import (
	"context"
	"time"
)

type UserData struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}
type UserStore interface {
	Create(ctx context.Context, createdData *UserData) error
	FindOne(ctx context.Context, filterBy *UserData, staging string) (*UserData, error)
}
