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
