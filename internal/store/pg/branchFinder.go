package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/lib/pq"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type BranchFinder struct {
	db *sql.DB
}

func NewBranchFinder(db *sql.DB) *BranchFinder {
	return &BranchFinder{db: db}
}

func (b *BranchFinder) FindMany(ctx context.Context, branchIDs []string) ([]store.BranchesData, error) {
	query := `
		SELECT 
			id, 
			tenant_id,
			name,
			address,
			created_at
		FROM branches
		WHERE id = ANY($1::uuid[])
	`
	rows, err := sqldb.WithinTxContextOrDB(ctx, b.db).QueryContext(ctx, query,
		pq.Array(branchIDs),
	)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return nil, err
	}
	branches := make([]store.BranchesData, 0)
	for rows.Next() {
		branch := store.BranchesData{}
		err = rows.Scan(
			&branch.BranchID,
			&branch.TenantID,
			&branch.Name,
			&branch.Address,
			&branch.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return branches, nil
}
