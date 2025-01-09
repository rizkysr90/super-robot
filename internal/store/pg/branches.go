package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type Branches struct {
	db *sql.DB
}

func NewBranches(db *sql.DB) *Branches {
	return &Branches{
		db: db,
	}
}

func (b *Branches) FindOne(ctx context.Context, filter *store.BranchesFilter) (*store.BranchesData, error) {
	query := `
		SELECT id, tenant_id, name, address, created_at 
FROM branches
WHERE (
    ($1 = '' OR id = $1::uuid) AND
    ($2 = '' OR name = $2) AND 
    ($3 = '' OR tenant_id = $3::uuid)
) AND deleted_at IS NULL
	`
	data := &store.BranchesData{}
	row := sqldb.WithinTxContextOrDB(ctx, b.db).
		QueryRowContext(ctx, query, filter.ID, filter.Name, filter.TenantID)
	if err := row.Err(); err != nil {
		return nil, err
	}
	err := row.Scan(&data.BranchID, &data.TenantID,
		&data.Name, &data.Address, &data.CreatedAt)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (b *Branches) Insert(ctx context.Context, branchData *store.BranchesData) error {
	query := `
		INSERT INTO branches (id, tenant_id, name, address, created_at, created_by, created_by_owner)
		VALUES (
			$1, 
			$2, 
			$3, 
			$4, 
			$5, 
			NULLIF($6, ''), 
			NULLIF($7, '')
		)
	`
	createFunc := func(tx sqldb.QueryExecutor) error {
		_, err := tx.ExecContext(ctx, query,
			branchData.BranchID,
			branchData.TenantID,
			branchData.Name,
			branchData.Address,
			branchData.CreatedAt,
			branchData.CreatedBy.String,
			branchData.CreatedByOwner.String,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return sqldb.WithinTxContextOrError(ctx, createFunc)
}

func (b *Branches) TotalBranches(ctx context.Context, tenantID string) (uint8, error) {
	result := 0
	query := `
		SELECT COUNT(id) FROM branches WHERE tenant_id = $1
		
	`
	row := sqldb.WithinTxContextOrDB(ctx, b.db).
		QueryRowContext(ctx, query, tenantID)
	if err := row.Err(); err != nil {
		return 0, err
	}
	err := row.Scan(&result)
	if err != nil {
		return 0, err
	}
	return uint8(result), nil
}
