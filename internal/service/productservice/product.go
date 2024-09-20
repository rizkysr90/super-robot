package productservice

import (
	"auth-service-rizkysr90-pos/internal/store"
	"database/sql"
)

type Service struct {
	db            *sql.DB
	productStore store.Product
}

func NewProductService(sqlDB *sql.DB, productStore store.Product) *Service {
	return &Service{
		db:            sqlDB,
		productStore: productStore,
	}
}