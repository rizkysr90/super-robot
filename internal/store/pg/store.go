package pg

import (
	"auth-service-rizkysr90-pos/internal/store"
	"context"
	"database/sql"
	"errors"
	"time"

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
			 contact, employee_id, created_at
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
	filter *store.StoreFilter, operationID string) (interface{}, error) {
	switch operationID {
	case "getallstore":
		const query = `
		WITH total_count AS (
			SELECT 
				COUNT(id) as total_elements
			FROM stores
			WHERE employee_id = $1 AND deleted_at IS NULL
		)
		SELECT 
			id, name, address, contact, created_at,
		(SELECT total_elements FROM total_count) AS total_elements
		FROM stores WHERE employee_id = $1 AND deleted_at IS NULL
		LIMIT $2
		OFFSET $3
	`
		var resultArr []store.StoreData
		rows, err := sqldb.WithinTxContextOrDB(ctx, s.db).QueryContext(ctx, query,
			filter.EmployeeID, filter.Pagination.PageSize, filter.Pagination.PageNumber)
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
	case "findbyname":
		const query = `
			SELECT name FROM stores WHERE name = $1 AND employee_id = $2 AND deleted_at IS NOT NULL
		`
		var result store.StoreData
		row := sqldb.WithinTxContextOrDB(ctx, s.db).QueryRowContext(ctx, query,
			filter.Name, filter.EmployeeID)
		err := row.Scan(&result.Name)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, errors.New("invalid operation id")
	}

}

func (s *Store) Delete(ctx context.Context, filter *store.StoreFilter) error {

	return sqldb.WithinTxContextOrError(ctx, func(tx sqldb.QueryExecutor) error {
		const query = `
			UPDATE stores SET deleted_at = $3 WHERE id = $1 AND employee_id = $2
		`
		_, err := tx.ExecContext(ctx, query,
			filter.StoreID,
			filter.EmployeeID,
			time.Now().UTC(),
		)
		if err != nil {
			return err
		}
		return nil
	})
}
func (s *Store) Update(ctx context.Context, updatedData *store.StoreData, filter *store.StoreFilter, operationID string) (int64, error) {
	var updateFunc func(tx sqldb.QueryExecutor) error
	var err error
	var rowsAffected int64
	switch operationID {
	case "updatestore":

		updateFunc = func(tx sqldb.QueryExecutor) error {
			query := `
				UPDATE stores SET 
					name = $3, 
					address = $4, 
					contact = $5
				WHERE id = $1 AND employee_id = $2 AND deleted_at IS NOT NULL
		`
			var result sql.Result
			result, err = tx.ExecContext(ctx,
				query,
				filter.StoreID,
				filter.EmployeeID,
				updatedData.Name,
				updatedData.Address,
				updatedData.Contact,
			)
			if err != nil {
				return err
			}
			rowsAffected, err = result.RowsAffected()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return rowsAffected, sqldb.WithinTxContextOrError(ctx, updateFunc)
}
