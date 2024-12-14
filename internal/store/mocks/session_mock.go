package mocks

import (
	"context"
	"rizkysr90-pos/internal/store"

	"github.com/stretchr/testify/mock"
)

// / MockSessionRedis is a mock implementation of the SessionRedis interface
type MockSessionRedis struct {
	mock.Mock
}

// Insert provides a mock function for inserting session data into Redis
func (m *MockSessionRedis) Insert(ctx context.Context, sessionData *store.SessionRedisData) error {
	args := m.Called(ctx, sessionData)
	return args.Error(0)
}

// Get provides a mock function for retrieving session data from Redis
func (m *MockSessionRedis) Get(ctx context.Context, sessionID string) (*store.SessionRedisData, error) {
	args := m.Called(ctx, sessionID)

	// Handle the case where first return value might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*store.SessionRedisData), args.Error(1)
}
