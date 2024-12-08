package store

import (
	"context"
	"database/sql"
)

type StateData struct {
	ID         string
	TenantName sql.NullString
}
type State interface {
	Insert(ctx context.Context, tenantData *StateData) error
	FindOne(ctx context.Context, id string) (*StateData, error)
	Delete(ctx context.Context, stateID string) error
}
