package store

import (
	"context"
	"time"
)

type ProductData struct {
	ProductID     string         
	ProductName   string        
	Price         float64        
	BasePrice     float64        
	StockQuantity int           
	CategoryID    string  		
	CreatedAt     time.Time      
	UpdatedAt    time.Time 
	DeletedAt    time.Time 
	Category *CategoryData
}
type FilterProduct struct {
	Limit      int
	Offset     int
	CategoryID string
}
type Product interface {
	Insert(ctx context.Context, insertedData *ProductData) error
	Update(ctx context.Context, updatedData *ProductData) error
	GetByName(ctx context.Context, productNameInput string) (*ProductData, error)
	GetAll(ctx context.Context, params *FilterProduct) ([]ProductData, int, error)
	GetByID(ctx context.Context, productIDInput string) (*ProductData, error)
	DeleteByID(ctx context.Context, productID string) error
}