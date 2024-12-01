package productservice

import (
	"database/sql"
	"rizkysr90-pos/internal/store"
	documentgen "rizkysr90-pos/pkg/documentGen"
)

type Service struct {
	db                *sql.DB
	productStore      store.Product
	documentGenerator documentgen.DocumentGenaratorInterface
}

func NewProductService(
	sqlDB *sql.DB,
	productStore store.Product,
	documentGenerator documentgen.DocumentGenaratorInterface,
) *Service {
	return &Service{
		db:                sqlDB,
		productStore:      productStore,
		documentGenerator: documentGenerator,
	}
}
