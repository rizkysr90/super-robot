package category

import (
	"database/sql"
	"rizkysr90-pos/internal/store"
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
