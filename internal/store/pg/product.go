package pg

import (
	"auth-service-rizkysr90-pos/internal/store"
	"context"
	"database/sql"
	"time"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)


type Product struct {
	db *sql.DB
}

func NewProduct(db *sql.DB) *Product {
	return &Product{
		db: db,
	}
}

func (p *Product) Insert(ctx context.Context, insertedData *store.ProductData) error {
	query := `
		INSERT INTO products (product_id, product_name, price, base_price, stock_quantity, category_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	createFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			insertedData.ProductID,
			insertedData.ProductName,
			insertedData.Price,
			insertedData.BasePrice,
			insertedData.StockQuantity,
			insertedData.CategoryID,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}

func (p *Product) Update(ctx context.Context, updatedData *store.ProductData) error {
	query := `
		UPDATE products
		SET 
			product_name = CASE WHEN $1::varchar IS NOT NULL THEN $1 ELSE product_name END,
			price = CASE WHEN $2::decimal IS NOT NULL THEN $2 ELSE price END,
			base_price = CASE WHEN $3::decimal IS NOT NULL THEN $3 ELSE base_price END,
			stock_quantity = CASE WHEN $4::int IS NOT NULL THEN $4 ELSE stock_quantity END,
			category_id = CASE WHEN $5::uuid IS NOT NULL THEN $5 ELSE category_id END,
			updated_at = CURRENT_TIMESTAMP
		WHERE 
			product_id = $6 AND deleted_at IS NULL;
	`
	updateFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			updatedData.ProductName,
			updatedData.Price,
			updatedData.BasePrice,
			updatedData.StockQuantity,
			updatedData.CategoryID,
			updatedData.ProductID,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, updateFunc)
}

func (p *Product) GetByName(ctx context.Context, productNameInput string) (*store.ProductData, error) {
	query := `
		SELECT product_id, product_name, price, base_price, stock_quantity,
		created_at, updated_at, deleted_at FROM products
		WHERE product_name = $1 AND deleted_at IS NULL
	`
	var (
		productID string
		productName string
		price float64
		basePrice float64
		stockQuantity int32
		createdAt time.Time
		updatedAt sql.NullTime
		deletedAt sql.NullTime
	)
	err := sqldb.WithinTxContextOrDB(ctx, p.db).QueryRowContext(ctx, query, productNameInput).Scan(
		&productID,
		&productName,
		&price,
		&basePrice,
		&stockQuantity,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &store.ProductData{
		ProductID: productID,
		ProductName: productName,
		Price: price,
		BasePrice: basePrice,
		StockQuantity: int(stockQuantity),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt.Time,
		DeletedAt: deletedAt.Time,
	}, nil
}
func (p *Product) GetAll(ctx context.Context, 
	params *store.FilterProduct) ([]store.ProductData, int, error) {
	query := `
		SELECT 
    p.product_id, 
    p.product_name, 
    p.price, 
    p.base_price, 
    p.stock_quantity, 
    p.category_id, 
    c.category_name, 
    p.created_at, 
    p.updated_at, 
    COUNT(p.product_id) OVER() AS total_count 
FROM 
    products p 
LEFT JOIN 
    categories c ON p.category_id = c.id 
WHERE 
    p.deleted_at IS NULL 
    AND (COALESCE($1, '') = '' OR p.category_id = $1::uuid)
ORDER BY 
    p.created_at DESC 
LIMIT $2 OFFSET $3
	`
	
	var products []store.ProductData
	var totalCount int

	rows, err := sqldb.WithinTxContextOrDB(ctx, p.db).QueryContext(
		ctx, query, params.CategoryID, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			productID     string
			productName   string
			price         float64
			basePrice     float64
			stockQuantity int
			categoryID    sql.NullString
			categoryName  sql.NullString
			createdAt     time.Time
			updatedAt     sql.NullTime
		)

		if err := rows.Scan(
			&productID,
			&productName,
			&price,
			&basePrice,
			&stockQuantity,
			&categoryID,
			&categoryName,
			&createdAt,
			&updatedAt,
			&totalCount,
		); err != nil {
			return nil, 0, err
		}

		product := store.ProductData{
			ProductID:     productID,
			ProductName:   productName,
			Price:         price,
			BasePrice:     basePrice,
			StockQuantity: stockQuantity,
			CreatedAt:     createdAt,
			UpdatedAt: updatedAt.Time,
			CategoryID: categoryID.String,
			Category: &store.CategoryData{
				ID: categoryID.String,
				CategoryName: categoryName.String,
			},
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return products, totalCount, nil
}

func (p *Product) GetByID(ctx context.Context, productIDInput string) (*store.ProductData, error) {
	query := `
		SELECT 
    p.product_id, 
    p.product_name, 
    p.price, 
    p.base_price, 
    p.stock_quantity, 
    p.category_id, 
    c.category_name, 
    p.created_at, 
    p.updated_at
FROM 
    products p 
LEFT JOIN 
    categories c ON p.category_id = c.id 
WHERE 
    p.deleted_at IS NULL 
    AND p.product_id = $1
	`
	var (
		productID     string
		productName   string
		price         float64
		basePrice     float64
		stockQuantity int
		categoryID    sql.NullString
		categoryName  sql.NullString
		createdAt     time.Time
		updatedAt     sql.NullTime
	)
	err := sqldb.WithinTxContextOrDB(ctx, p.db).QueryRowContext(ctx, query, productIDInput).Scan(
		&productID,
			&productName,
			&price,
			&basePrice,
			&stockQuantity,
			&categoryID,
			&categoryName,
			&createdAt,
			&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	product := store.ProductData{
		ProductID:     productID,
		ProductName:   productName,
		Price:         price,
		BasePrice:     basePrice,
		StockQuantity: stockQuantity,
		CreatedAt:     createdAt,
		UpdatedAt: updatedAt.Time,
		CategoryID: categoryID.String,
		Category: &store.CategoryData{
			ID: categoryID.String,
			CategoryName: categoryName.String,
		},
	}
	return &product, nil
}
func (p *Product) DeleteByID(ctx context.Context, productID string) error {
	query := `
		UPDATE products SET deleted_at = CURRENT_TIMESTAMP WHERE product_id = $1 AND deleted_at IS NULL
	`
	updateFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			productID,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, updateFunc)
	
}

