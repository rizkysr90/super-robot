package mocks

import (
	"context"
	"rizkysr90-pos/internal/store"

	"github.com/stretchr/testify/mock"
)

// MockTenant is a mock implementation of the Tenant interface
type MockTenant struct {
	mock.Mock
}

// Insert provides a mock function for inserting tenant data
func (m *MockTenant) Insert(ctx context.Context, tenantData *store.TenantData) error {
	args := m.Called(ctx, tenantData)
	return args.Error(0)
}

// Update provides a mock function for updating tenant data
func (m *MockTenant) Update(ctx context.Context, tenantData *store.TenantData) error {
	args := m.Called(ctx, tenantData)
	return args.Error(0)
}
