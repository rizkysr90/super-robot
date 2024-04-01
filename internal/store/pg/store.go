package pg

import (
	"auth-service-rizkysr90-pos/internal/store"
	"context"
	"database/sql"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type Store struct {
	db *sql.DB
}

func NewStoreDB(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Insert(ctx context.Context, data *store.StoreData) error {
	createFunc := func(tx sqldb.QueryExecutor) error {
		query := `
			INSERT INTO stores 
			(id, name, address,
			 contact, user_id, created_at
			)
			VALUES 
			($1, $2, $3, $4, $5, $6)
		`
		_, err := tx.ExecContext(ctx, query,
			data.ID,
			data.Name,
			data.Address,
			data.Contact,
			data.UserID,
			data.CreatedAt,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}

func (s *Store) Finder(ctx context.Context,
	filter *store.StoreFilter, operationID string) ([]store.StoreData, error) {
	const query = `
		WITH total_count AS (
			SELECT 
				COUNT(s.id) as total_elements
			FROM stores s JOIN employees e ON s.id = e.store_id 
			WHERE e.id = $1
		)
		SELECT 
			s.id, s.name, s.address, s.contact, s.created_at,
		(SELECT total_elements FROM total_count) AS total_elements
		FROM stores s JOIN employees e ON s.id = e.store_id 
		WHERE e.id = $1
		LIMIT $2
		OFFSET $3
	`
	var resultArr []store.StoreData
	rows, err := sqldb.WithinTxContextOrDB(ctx, s.db).QueryContext(ctx, query,
		filter.UserID, filter.Pagination.PageSize, filter.Pagination.PageNumber)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result store.StoreData
		err = rows.Scan(&result.ID, &result.Name, &result.Address,
			&result.Contact, &result.CreatedAt, &filter.Pagination.TotalElements)
		if err != nil {
			return nil, err
		}
		resultArr = append(resultArr, result)
	}
	return resultArr, nil
}
