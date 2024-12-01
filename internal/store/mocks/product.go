package mocks

import (
	"context"

	"rizkysr90-pos/internal/store"

	"github.com/stretchr/testify/mock"
)

// MockProductStore is a mock implementation of Product interface
type MockProductStore struct {
	mock.Mock
}

// Insert provides a mock function with given fields: ctx, insertedData
func (m *MockProductStore) Insert(ctx context.Context, insertedData *store.ProductData) error {
	ret := m.Called(ctx, insertedData)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *store.ProductData) error); ok {
		r0 = rf(ctx, insertedData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, updatedData
func (m *MockProductStore) Update(ctx context.Context, updatedData *store.ProductData) error {
	ret := m.Called(ctx, updatedData)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *store.ProductData) error); ok {
		r0 = rf(ctx, updatedData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByName provides a mock function with given fields: ctx, productNameInput
func (m *MockProductStore) GetByName(ctx context.Context, productNameInput string) (*store.ProductData, error) {
	ret := m.Called(ctx, productNameInput)

	var r0 *store.ProductData
	if rf, ok := ret.Get(0).(func(context.Context, string) *store.ProductData); ok {
		r0 = rf(ctx, productNameInput)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*store.ProductData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, productNameInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx, params
func (m *MockProductStore) GetAll(ctx context.Context, params *store.FilterProduct) ([]store.ProductData, int, error) {
	ret := m.Called(ctx, params)

	var r0 []store.ProductData
	if rf, ok := ret.Get(0).(func(context.Context, *store.FilterProduct) []store.ProductData); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]store.ProductData)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(context.Context, *store.FilterProduct) int); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *store.FilterProduct) error); ok {
		r2 = rf(ctx, params)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetByID provides a mock function with given fields: ctx, productIDInput
func (m *MockProductStore) GetByID(ctx context.Context, productIDInput string) (*store.ProductData, error) {
	ret := m.Called(ctx, productIDInput)

	var r0 *store.ProductData
	if rf, ok := ret.Get(0).(func(context.Context, string) *store.ProductData); ok {
		r0 = rf(ctx, productIDInput)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*store.ProductData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, productIDInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByID provides a mock function with given fields: ctx, productID
func (m *MockProductStore) DeleteByID(ctx context.Context, productID string) error {
	ret := m.Called(ctx, productID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, productID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
