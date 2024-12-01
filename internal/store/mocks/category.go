//nolint:forcetypeassert
package mocks

import (
	"context"

	"rizkysr90-pos/internal/store"

	"github.com/stretchr/testify/mock"
)

type MockCategoryStore struct {
	mock.Mock
}

func (m *MockCategoryStore) Create(ctx context.Context, category *store.CategoryData) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryStore) Update(ctx context.Context, category *store.CategoryData) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryStore) SoftDelete(ctx context.Context, categoryID string) error {
	args := m.Called(ctx, categoryID)
	return args.Error(0)
}

func (m *MockCategoryStore) FindByID(ctx context.Context, categoryID string) (*store.CategoryData, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) != nil {
		return args.Get(0).(*store.CategoryData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCategoryStore) FindAllPagination(ctx context.Context, p *store.Pagination) ([]store.CategoryData, error) {
	args := m.Called(ctx, p)
	if args.Get(0) != nil {
		return args.Get(0).([]store.CategoryData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCategoryStore) FindByName(ctx context.Context, categoryName string) (*store.CategoryData, error) {
	args := m.Called(ctx, categoryName)
	if args.Get(0) != nil {
		return args.Get(0).(*store.CategoryData), args.Error(1)
	}
	return nil, args.Error(1)
}
