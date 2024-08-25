package pg

import (
	"auth-service-rizkysr90-pos/internal/store"
	"context"
	"database/sql"
	"time"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type Category struct {
	db *sql.DB
}

func NewCategory(db *sql.DB) *Category {
	return &Category{
		db: db,
	}
}

func (c *Category) Create(ctx context.Context, category *store.CategoryData) error {
	query := `
		INSERT INTO categories (id, category_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	createFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			category.ID,
			category.CategoryName,
			category.CreatedAt,
			category.UpdatedAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}
func (c *Category) Update(ctx context.Context, category *store.CategoryData) error {
	query := `
		UPDATE categories SET category_name = $2, updated_at = $3 WHERE id = $1 AND deleted_at IS NULL
	`
	updateFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			category.ID,
			category.CategoryName,
			category.UpdatedAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, updateFunc)
}
func (c *Category) SoftDelete(ctx context.Context, categoryID string) error {
	query := `
		UPDATE categories SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL;
	`
	softDeleteFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			categoryID,
			time.Now().UTC(),
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, softDeleteFunc)
}
func (c *Category) FindByID(ctx context.Context, categoryID string) (*store.CategoryData, error) {
	query := `
		SELECT id, category_name, created_at, updated_at, deleted_at FROM categories
		WHERE id = $1 AND deleted_at IS NULL
	`
	data := &store.CategoryData{}
	row := sqldb.WithinTxContextOrDB(ctx, c.db).
		QueryRowContext(ctx, query, categoryID)
	if err := row.Err(); err != nil {
		return nil, err
	}
	err := row.Scan(&data.ID, &data.CategoryName,
		&data.CreatedAt, &data.UpdatedAt, &data.DeletedAt)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (c *Category) FindAllPagination(ctx context.Context, p *store.Pagination) ([]store.CategoryData, error) {
	query := `
		WITH total AS (
			SELECT COUNT(id) AS total_elements
			FROM categories
			WHERE deleted_at IS NULL
		)
		SELECT c.id, c.category_name, c.created_at, c.updated_at, c.deleted_at, t.total_elements
		FROM categories c, total t
		WHERE c.deleted_at IS NULL
		ORDER BY c.category_name DESC
		LIMIT $1 OFFSET $2;
	`
	// default page number is 1
	offset := (p.PageNumber - 1) * p.PageSize
	rows, err := sqldb.WithinTxContextOrDB(ctx, c.db).QueryContext(ctx, query,
		p.PageSize,
		offset,
	)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return nil, err
	}
	categories := make([]store.CategoryData, 0)
	for rows.Next() {
		category := store.CategoryData{Pagination: &store.Pagination{}}
		var totalElement int
		err = rows.Scan(
			&category.ID,
			&category.CategoryName,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
			&totalElement,
		)
		category.Pagination.TotalElements = totalElement

		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *Category) FindByName(ctx context.Context, categoryName string) (*store.CategoryData, error) {
	query := `
		SELECT id, category_name, created_at, updated_at, deleted_at FROM categories
		WHERE category_name = $1 AND deleted_at IS NULL
	`
	data := &store.CategoryData{}
	row := sqldb.WithinTxContextOrDB(ctx, c.db).
		QueryRowContext(ctx, query, categoryName)
	if err := row.Err(); err != nil {
		return nil, err
	}
	err := row.Scan(&data.ID, &data.CategoryName,
		&data.CreatedAt, &data.UpdatedAt, &data.DeletedAt)
	if err != nil {
		return nil, err
	}
	return data, nil
}
