package mocks

import (
	"rizkysr90-pos/internal/store"

	"github.com/stretchr/testify/mock"
)

// MockDocumentGenerator is a mock implementation of DocumentGeneratorInterface
type MockDocumentGenerator struct {
	mock.Mock
}

// LabelPricing mocks the label pricing document generation
func (m *MockDocumentGenerator) LabelPricing(productData *store.ProductData) ([]byte, error) {
	args := m.Called(productData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
