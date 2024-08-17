package store

import (
	"context"
	"database/sql"
	"time"
)

type CategoryData struct {
	Id string
	CategoryName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime

	Pagination *Pagination
}
type Category interface {
	Create(ctx context.Context, category *CategoryData) error
	Update(ctx context.Context, category *CategoryData) error
	SoftDelete(ctx context.Context, categoryId string) error
	FindById(ctx context.Context, categoryId string) (*CategoryData, error)
	FindAllPagination(ctx context.Context, p *Pagination) ([]CategoryData, error)
	FindByName(ctx context.Context, categoryName string) (*CategoryData, error)
}