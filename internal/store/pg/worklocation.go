package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type WorkLocation struct {
	db *sql.DB
}

func NewWorkLocation(db *sql.DB) *WorkLocation {
	return &WorkLocation{
		db: db,
	}
}
func (w *WorkLocation) FindByUserID(ctx context.Context, userID string) ([]store.WorkLocationData, error) {
	query := `
		SELECT 
			id, tenant_id, user_id, branch_id, 
			created_at, updated_at, deleted_at,
			created_by, updated_by, deleted_by
		FROM user_assignments
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	rows, err := sqldb.WithinTxContextOrDB(ctx, w.db).QueryContext(ctx, query,
		userID,
	)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return nil, err
	}
	workLocations := make([]store.WorkLocationData, 0)
	for rows.Next() {
		data := store.WorkLocationData{}
		err = rows.Scan(
			&data.ID,
			&data.TenantID,
			&data.UserID,
			&data.BranchID,
			&data.CreatedAt,
			&data.UpdatedAt,
			&data.DeletedAt,
			&data.CreatedBy,
			&data.UpdatedBy,
			&data.DeletedBy,
		)
		if err != nil {
			return nil, err
		}
		workLocations = append(workLocations, data)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return workLocations, nil
}
