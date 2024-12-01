package store

import (
	"context"
	"time"
)

type ProductData struct {
	Category      *CategoryData
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	ProductID     string
	ProductName   string
	CategoryID    string
	Price         float64
	BasePrice     float64
	StockQuantity int
}
type FilterProduct struct {
	CategoryID string
	Limit      int
	Offset     int
}
type Product interface {
	Insert(ctx context.Context, insertedData *ProductData) error
	Update(ctx context.Context, updatedData *ProductData) error
	GetByName(ctx context.Context, productNameInput string) (*ProductData, error)
	GetAll(ctx context.Context, params *FilterProduct) ([]ProductData, int, error)
	GetByID(ctx context.Context, productIDInput string) (*ProductData, error)
	DeleteByID(ctx context.Context, productID string) error
}
