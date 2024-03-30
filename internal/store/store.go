package store

import (
	"context"
	"database/sql"
	"time"
)

type StoreData struct {
	ID        string
	Name      string
	Address   string
	Contact   string
	UserID    string
	CreatedAt time.Time
	DeletedAt sql.NullTime
}

type StoreStore interface {
	Insert(ctx context.Context, data *StoreData) error
}
