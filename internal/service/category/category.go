package category

import (
	"auth-service-rizkysr90-pos/internal/store"
	"database/sql"
)

type CategoryService struct {
	db *sql.DB
	categoryStore store.Category
}
func NewCategoryService(sqlDB *sql.DB, categoryStore store.Category) *CategoryService {
	return &CategoryService{
		db: sqlDB,
		categoryStore: categoryStore,
	}
}
