package mocks

import (
	"context"
	"rizkysr90-pos/internal/store"

	"github.com/stretchr/testify/mock"
)

// MockBranchStore represents a mock implementation of store.Branch interface
type MockBranchStore struct {
	mock.Mock
}

func (m *MockBranchStore) FindOne(ctx context.Context, filter *store.BranchesFilter) (*store.BranchesData, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.BranchesData), args.Error(1)
}

func (m *MockBranchStore) Insert(ctx context.Context, data *store.BranchesData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}
func (m *MockBranchStore) TotalBranches(ctx context.Context, tenantID string) (uint8, error) {
	args := m.Called(ctx, tenantID)
	return uint8(args.Int(0)), args.Error(1)
}

// FindManyWithPaginated mocks the FindManyWithPaginated method
func (m *MockBranchStore) FindManyWithPaginated(ctx context.Context,
	filter *store.BranchesFilter) ([]store.BranchesData, *store.Pagination, error) {
	args := m.Called(ctx, filter)

	// Handle the return values
	var branches []store.BranchesData
	var pagination *store.Pagination

	// Type assert the first return value if it's not nil
	if args.Get(0) != nil {
		branches = args.Get(0).([]store.BranchesData)
	}

	// Type assert the second return value if it's not nil
	if args.Get(1) != nil {
		pagination = args.Get(1).(*store.Pagination)
	}

	return branches, pagination, args.Error(2)
}
