package mocks

import (
	"context"
	"rizkysr90-pos/internal/store"

	"github.com/stretchr/testify/mock"
)

// MockStateStore is a mock implementation of the State interface
type MockStateStore struct {
	mock.Mock
}

// Insert mocks the Insert method
func (m *MockStateStore) Insert(ctx context.Context, stateData *store.StateData) error {
	args := m.Called(ctx, stateData)
	return args.Error(0)
}

// FindOne mocks the FindOne method
func (m *MockStateStore) FindOne(ctx context.Context, id string) (*store.StateData, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.StateData), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockStateStore) Delete(ctx context.Context, stateID string) error {
	args := m.Called(ctx, stateID)
	return args.Error(0)
}
