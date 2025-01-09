package mocks

import (
	"context"
	"rizkysr90-pos/internal/store"

	"github.com/stretchr/testify/mock"
)

// MockUser is a mock implementation of the User interface
type MockUser struct {
	mock.Mock
}

// Insert provides a mock function for inserting user data
func (m *MockUser) Insert(ctx context.Context, userData *store.UserData) error {
	args := m.Called(ctx, userData)
	return args.Error(0)
}

func (m *MockUser) FindOne(ctx context.Context, filter *store.UserQueryFilter) (*store.UserData, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.UserData), args.Error(1)
}

// Update provides a mock function for updating user data
func (m *MockUser) Update(ctx context.Context, userData *store.UserData) error {
	args := m.Called(ctx, userData)
	return args.Error(0)
}
