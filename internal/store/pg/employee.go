package pg

import (
	"auth-service-rizkysr90-pos/internal/store"
	"context"
	"database/sql"
	"errors"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type Employee struct {
	db *sql.DB
}

func NewEmployeeDB(db *sql.DB) *Employee {
	return &Employee{
		db: db,
	}
}

func (e *Employee) Insert(ctx context.Context, data *store.EmployeeData) error {
	createFunc := func(tx sqldb.QueryExecutor) error {
		var err error
		if data.UserID != "" {
			query := `
			INSERT INTO employees 
			(id, name, contact,
			 username, password, store_id, 
			 role, created_at, user_id
			)
			VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`
			_, err = tx.ExecContext(ctx, query,
				data.ID,
				data.Name,
				data.Contact,
				data.Username,
				data.Password,
				data.StoreID,
				data.Role,
				data.CreatedAt,
				data.UserID,
			)
		} else {
			query := `
			INSERT INTO employees 
			(id, name, contact,
			 username, password, store_id, 
			 role, created_at
			)
			VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		`
			_, err = tx.ExecContext(ctx, query,
				data.ID,
				data.Name,
				data.Contact,
				data.Username,
				data.Password,
				data.StoreID,
				data.Role,
				data.CreatedAt,
			)
		}
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}

func (e *Employee) Finder(ctx context.Context,
	filter *store.EmployeeFilter, operationID string) ([]store.EmployeeData, error) {
	// const query = `
	// 	WITH total_count AS (
	// 		SELECT COUNT(id) AS total_elements FROM stores WHERE user_id = $1
	// 	)
	// 	SELECT
	// 	id, name, address, contact, created_at,
	// 	(SELECT total_elements FROM total_count) AS total_elements
	// 	FROM stores WHERE user_id = $1
	// 	LIMIT $2
	// 	OFFSET $3
	// `
	// var resultArr []store.StoreData
	// rows, err := sqldb.WithinTxContextOrDB(ctx, s.db).QueryContext(ctx, query,
	// 	filter.UserID, filter.Pagination.PageSize, filter.Pagination.PageNumber)
	// if err != nil {
	// 	return nil, err
	// }
	// if err = rows.Err(); err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var result store.StoreData
	// 	err = rows.Scan(&result.ID, &result.Name, &result.Address,
	// 		&result.Contact, &result.CreatedAt, &filter.Pagination.TotalElements)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	resultArr = append(resultArr, result)
	// }
	// return resultArr, nil
	return nil, nil
}
func (e *Employee) FindOne(ctx context.Context,
	filterBy *store.EmployeeFilter, operationID string) (*store.EmployeeData, error) {
	var result store.EmployeeData
	switch operationID {
	case "findactiveemployee":
		// 1 is for filter by email
		query := `SELECT id, password, role
				 FROM employees WHERE username = $1`
		err := sqldb.WithinTxContextOrDB(ctx, e.db).
			QueryRowContext(ctx, query, filterBy.Username).
			Scan(&result.ID, &result.Password, &result.Role)
		if err != nil {
			return nil, err
		}
		return &result, nil
	case "findByID":
		query := `SELECT id, password, role, refresh_token
				 FROM employees WHERE id = $1`
		err := sqldb.WithinTxContextOrDB(ctx, e.db).
			QueryRowContext(ctx, query, filterBy.ID).
			Scan(&result.ID, &result.Password, &result.Role, &result.RefreshToken)
		if err != nil {
			return nil, err
		}
		return &result, nil
	default:
		return nil, errors.New("staging db not found")
	}

}
func updateRefreshtokenEmployee(ctx context.Context,
	updatedData *store.EmployeeData, f *store.EmployeeFilter) func(tx sqldb.QueryExecutor) error {

	result := func(tx sqldb.QueryExecutor) error {
		query := `
			UPDATE employees SET refresh_token = $1
			WHERE username = $2
		`
		_, err := tx.ExecContext(ctx, query, updatedData.RefreshToken, f.Username)
		if err != nil {
			return err
		}
		return nil
	}
	return result
}
func (u *Employee) Update(ctx context.Context,
	updatedData *store.EmployeeData,
	filter *store.EmployeeFilter,
	staging string) error {

	var updateFunc func(tx sqldb.QueryExecutor) error
	switch staging {
	case "updaterefreshtoken":
		updateFunc = updateRefreshtokenEmployee(ctx, updatedData, filter)
	default:
		return errors.New("staging db not found")
	}
	return sqldb.WithinTxContextOrError(ctx, updateFunc)
}
