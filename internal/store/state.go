package store

import "context"

type StateData struct {
	ID string
}
type State interface {
	Insert(ctx context.Context, id string) error
	FindOne(ctx context.Context, id string) (*StateData, error)
}
