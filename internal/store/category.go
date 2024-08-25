package store

import (
	"context"
	"database/sql"
	"time"
)

type CategoryData struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Pagination   *Pagination
	DeletedAt    sql.NullTime
	ID           string
	CategoryName string
}

type Category interface {
	Create(ctx context.Context, category *CategoryData) error
	Update(ctx context.Context, category *CategoryData) error
	SoftDelete(ctx context.Context, categoryID string) error
	FindByID(ctx context.Context, categoryID string) (*CategoryData, error)
	FindAllPagination(ctx context.Context, p *Pagination) ([]CategoryData, error)
	FindByName(ctx context.Context, categoryName string) (*CategoryData, error)
}
