package pg

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"

	"github.com/lib/pq"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type AssignmentRoles struct {
	db *sql.DB
}

func NewAssignmentRoles(db *sql.DB) *AssignmentRoles {
	return &AssignmentRoles{
		db: db,
	}
}

func (a *AssignmentRoles) FindMany(ctx context.Context,
	assignmentIDs []string) ([]store.AssignmentRoleData, error) {

	query := `
		SELECT id, tenant_id, assignment_id, tenant_role_id,
			created_at, updated_at, deleted_at,
			created_by, updated_by, deleted_by
		FROM user_assignment_roles
		WHERE assignment_id = ANY($1::UUID[]) AND deleted_at IS NULL
	`
	rows, err := sqldb.WithinTxContextOrDB(ctx, a.db).QueryContext(ctx, query,
		pq.Array(assignmentIDs),
	)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return nil, err
	}
	roles := make([]store.AssignmentRoleData, 0)
	for rows.Next() {
		data := store.AssignmentRoleData{}
		err = rows.Scan(
			&data.ID,
			&data.TenantID,
			&data.AssignmentID,
			&data.TenantRoleID,
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
		roles = append(roles, data)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
