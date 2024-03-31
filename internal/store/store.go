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
type SetResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Contact   string    `json:"contact"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
type StoreFilter struct {
	UserID     string
	Pagination Pagination
}

func (s *StoreData) ToPayloadResponse(data *StoreData) *SetResponse {
	return &SetResponse{
		ID:        s.ID,
		Name:      s.Name,
		Address:   s.Address,
		Contact:   s.Contact,
		CreatedAt: s.CreatedAt,
		DeletedAt: s.DeletedAt.Time,
	}
}

type StoreStore interface {
	Insert(ctx context.Context, data *StoreData) error
	Finder(ctx context.Context, filter *StoreFilter, operationID string) ([]StoreData, error)
}
