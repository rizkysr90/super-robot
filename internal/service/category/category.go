package category

import (
	"auth-service-rizkysr90-pos/internal/store"
	"database/sql"
)

type Service struct {
	db            *sql.DB
	categoryStore store.Category
}

func NewCategoryService(sqlDB *sql.DB, categoryStore store.Category) *Service {
	return &Service{
		db:            sqlDB,
		categoryStore: categoryStore,
	}
}
