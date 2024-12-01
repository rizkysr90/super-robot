package payload

import "time"

type GenerateBarcodeRequest struct {
	ProductID string `json:"product_id"`
}
type GenerateBarcodeResponse struct {
	PDFBytes []byte `json:"pdf_bytes"`
}

type CategoryData struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
	ID           string    `json:"id"`
	CategoryName string    `json:"category_name"`
}

type ReqCreateCategory struct {
	CategoryName string `json:"category_name"`
}

type ResCreateCategory struct{}

type ReqUpdateCategory struct {
	ID           string `json:"id"`
	CategoryName string `json:"category_name"`
}

type ResUpdateCategory struct{}

type ReqGetCategoryByID struct {
	CategoryID string `json:"id"`
}

type ResGetCategoryByID struct {
	*CategoryData `json:"data"`
}

// ReqGetAllCategory represents the request payload for getting all categories.
type ReqGetAllCategory struct {
	PageSize   int `example:"20" json:"page_size"`
	PageNumber int `example:"1"  json:"page_number"`
}

// ResGetAllCategory represents the response payload for getting all categories.
type ResGetAllCategory struct {
	Data     []CategoryData `json:"data"`
	Metadata Pagination     `json:"metadata"`
}

type ReqDeleteCategory struct {
	ID string `json:"id"`
}

type ResDeleteCategory struct{}

type Pagination struct {
	PageSize      int `json:"page_size"`
	PageNumber    int `json:"page_number"`
	TotalPages    int `json:"total_pages"`
	TotalElements int `json:"total_elements"`
}
type ProductData struct {
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
	ProductID     string    `json:"product_id"`
	ProductName   string    `json:"product_name"`
	CategoryID    string    `json:"category_id"`
	CategoryName  string    `json:"category_name"`
	Price         float64   `json:"price"`
	BasePrice     float64   `json:"base_price"`
	StockQuantity int       `json:"stock_quantity"`
}

type ReqCreateProduct struct {
	ProductName   string  `json:"product_name"   validate:"required,max=40"`
	CategoryID    string  `json:"category_id"    validate:"required"`
	Price         float64 `json:"price"          validate:"gte=0"`
	BasePrice     float64 `json:"base_price"     validate:"gte=0"`
	StockQuantity int     `json:"stock_quantity" validate:"gte=0"`
}

type ResCreateProduct struct{}

type ReqUpdateProduct struct {
	ProductID     string  `json:"product_id"     validate:"required"`
	ProductName   string  `json:"product_name"   validate:"required,max=40"`
	CategoryID    string  `json:"category_id"    validate:"required,uuid"`
	Price         float64 `json:"price"          validate:"gte=0"`
	BasePrice     float64 `json:"base_price"     validate:"gte=0"`
	StockQuantity int     `json:"stock_quantity" validate:"gte=0"`
}

type ResUpdateProduct struct{}

type ReqGetAllProducts struct {
	CategoryID string `json:"category_id,omitempty"`
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
}

type ResGetAllProducts struct {
	Data     []ProductData `json:"data"`
	Metadata Pagination    `json:"metadata"`
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

type ResDeleteProductByID struct{}
