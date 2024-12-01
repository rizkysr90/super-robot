package productservice

import (
	"database/sql"
	"rizkysr90-pos/internal/store"
)

type Service struct {
	db           *sql.DB
	productStore store.Product
}

func NewProductService(sqlDB *sql.DB, productStore store.Product) *Service {
	return &Service{
		db:           sqlDB,
		productStore: productStore,
	}
}
