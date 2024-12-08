package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type State struct {
	db *sql.DB
}

func NewState(db *sql.DB) *State {
	return &State{
		db: db,
	}
}
func (s *State) Insert(ctx context.Context, stateData *store.StateData) error {
	query := `
		INSERT INTO states (state_id, tenant_name)
			VALUES ($1, NULLIF($2, ''));
	`
	createFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			stateData.ID,
			stateData.TenantName,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}
func (s *State) FindOne(ctx context.Context, stateID string) (*store.StateData, error) {
	query := `
		SELECT state_id, tenant_name FROM states
		WHERE state_id = $1
	`
	data := &store.StateData{}
	row := sqldb.WithinTxContextOrDB(ctx, s.db).
		QueryRowContext(ctx, query, stateID)
	if err := row.Err(); err != nil {
		return nil, err
	}
	err := row.Scan(&data.ID, &data.TenantName)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (s *State) Delete(ctx context.Context, stateID string) error {
	query := `
		DELETE FROM states WHERE state_id = $1
	`
	deleteFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			stateID,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, deleteFunc)
}
