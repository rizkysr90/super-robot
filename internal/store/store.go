package store

import (
	"context"
	"database/sql"
	"html"
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
	StoreID    string
	EmployeeID string
	Name       string
	Pagination Pagination
}

func (s *StoreData) ToPayloadResponse(data *StoreData) *SetResponse {
	return &SetResponse{
		ID:        s.ID,
		Name:      html.UnescapeString(s.Name),
		Address:   html.UnescapeString(s.Address),
		Contact:   html.UnescapeString(s.Contact),
		CreatedAt: s.CreatedAt,
		DeletedAt: s.DeletedAt.Time,
	}
}

type StoreStore interface {
	Insert(ctx context.Context, data *StoreData) error
	Finder(ctx context.Context, filter *StoreFilter, operationID string) (interface{}, error)
	Delete(ctx context.Context, filter *StoreFilter) error
	Update(ctx context.Context, updatedData *StoreData, filter *StoreFilter, operationID string) (int64, error)
}
