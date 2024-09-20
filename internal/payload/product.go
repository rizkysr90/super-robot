package payload

import "time"

type ProductData struct {
	ProductID     string    `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Price         float64   `json:"price"`
	BasePrice     float64   `json:"base_price"`
	StockQuantity int       `json:"stock_quantity"`
	CategoryID    string    `json:"category_id"`
	CategoryName  string    `json:"category_name"` // From categories table
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}
type ReqCreateProduct struct {
	ProductName   string  `json:"product_name" validate:"required,max=40"`
	Price         float64 `json:"price" validate:"gte=0"`
	BasePrice     float64 `json:"base_price" validate:"gte=0"`
	StockQuantity int     `json:"stock_quantity" validate:"gte=0"`
	CategoryID    string  `json:"category_id" validate:"required"`
}
type ResCreateProduct struct {
	
}
type ReqUpdateProduct struct {
	ProductID     string  `json:"product_id" validate:"required"`
	ProductName   string  `json:"product_name" validate:"required,max=40"`
	Price         float64 `json:"price" validate:"gte=0"`
	BasePrice     float64 `json:"base_price" validate:"gte=0"`
	StockQuantity int     `json:"stock_quantity" validate:"gte=0"`
	CategoryID    string  `json:"category_id" validate:"required,uuid"`
}

type ResUpdateProduct struct {

}

type ReqGetAllProducts struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	CategoryID string `json:"category_id,omitempty"` // Optional category filter
}

type ResGetAllProducts struct {
	Data     []ProductData `json:"data"`
	Metadata Pagination     `json:"metadata"`
}

type ResGetProductByID struct {
	Data ProductData `json:"data"`
}

type ReqGetProductByID struct {
	ProductID string `json:"product_id" validate:"required"`
}
type ReqDeleteProductByID struct {
	ProductID string `json:"product_id" validate:"required"`
}

type ResDeleteProductByID struct {
}