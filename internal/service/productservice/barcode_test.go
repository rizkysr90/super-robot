package productservice

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/store/mocks"
	"rizkysr90-pos/pkg/errorHandler"
)

func TestService_GenerateBarcodePDF(t *testing.T) {
	// Setup test data
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	samplePDFBytes := []byte("sample pdf content")

	tests := []struct {
		name             string
		input            *payload.GenerateBarcodeRequest
		setupMocks       func(mockProductStore *mocks.MockProductStore, mockDocGen *mocks.MockDocumentGenerator)
		expectedError    error
		validateResponse func(t *testing.T, response *payload.GenerateBarcodeResponse)
	}{
		{
			name: "Success - Valid Product ID",
			input: &payload.GenerateBarcodeRequest{
				ProductID: "PROD123",
			},
			setupMocks: func(mockProductStore *mocks.MockProductStore, mockDocGen *mocks.MockDocumentGenerator) {
				productData := &store.ProductData{
					ProductID:     "PROD123",
					ProductName:   "Test Product",
					Price:         100.0,
					BasePrice:     80.0,
					StockQuantity: 50,
					CreatedAt:     fixedTime,
					UpdatedAt:     fixedTime,
				}

				mockProductStore.On("GetByID", mock.Anything, "PROD123").
					Return(productData, nil)

				mockDocGen.On("LabelPricing", productData).
					Return(samplePDFBytes, nil)
			},
			expectedError: nil,
			validateResponse: func(t *testing.T, response *payload.GenerateBarcodeResponse) {
				assert.NotNil(t, response)
				assert.Equal(t, samplePDFBytes, response.PDFBytes)
			},
		},
		{
			name: "Error - Product Not Found",
			input: &payload.GenerateBarcodeRequest{
				ProductID: "NONEXISTENT",
			},
			setupMocks: func(mockProductStore *mocks.MockProductStore, mockDocGen *mocks.MockDocumentGenerator) {
				mockProductStore.On("GetByID", mock.Anything, "NONEXISTENT").
					Return(nil, sql.ErrNoRows)
			},
			expectedError: errorHandler.NewInternalServer(
				errorHandler.WithInfo("failed to fetch product"),
				errorHandler.WithMessage(sql.ErrNoRows.Error()),
			),
			validateResponse: func(t *testing.T, response *payload.GenerateBarcodeResponse) {
				assert.Nil(t, response)
			},
		},
		{
			name: "Error - PDF Generation Failed",
			input: &payload.GenerateBarcodeRequest{
				ProductID: "PROD123",
			},
			setupMocks: func(mockProductStore *mocks.MockProductStore, mockDocGen *mocks.MockDocumentGenerator) {
				productData := &store.ProductData{
					ProductID:     "PROD123",
					ProductName:   "Test Product",
					Price:         100.0,
					BasePrice:     80.0,
					StockQuantity: 50,
				}

				mockProductStore.On("GetByID", mock.Anything, "PROD123").
					Return(productData, nil)

				mockDocGen.On("LabelPricing", productData).
					Return(nil, assert.AnError)
			},
			expectedError: errorHandler.NewInternalServer(
				errorHandler.WithInfo("failed to generate barcode product"),
				errorHandler.WithMessage(assert.AnError.Error()),
			),
			validateResponse: func(t *testing.T, response *payload.GenerateBarcodeResponse) {
				assert.Nil(t, response)
			},
		},
		{
			name: "Error - Database Connection Error",
			input: &payload.GenerateBarcodeRequest{
				ProductID: "PROD123",
			},
			setupMocks: func(mockProductStore *mocks.MockProductStore, mockDocGen *mocks.MockDocumentGenerator) {
				mockProductStore.On("GetByID", mock.Anything, "PROD123").
					Return(nil, sql.ErrConnDone)
			},
			expectedError: errorHandler.NewInternalServer(
				errorHandler.WithInfo("failed to fetch product"),
				errorHandler.WithMessage(sql.ErrConnDone.Error()),
			),
			validateResponse: func(t *testing.T, response *payload.GenerateBarcodeResponse) {
				assert.Nil(t, response)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockProductStore := new(mocks.MockProductStore)
			mockDocGen := new(mocks.MockDocumentGenerator)

			tt.setupMocks(mockProductStore, mockDocGen)

			// Create service with mocked dependencies
			svc := &Service{
				productStore:      mockProductStore,
				documentGenerator: mockDocGen,
			}

			// Execute
			result, err := svc.GenerateBarcodePDF(context.Background(), tt.input)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Validate response
			tt.validateResponse(t, result)

			// Verify all mock expectations were met
			mockProductStore.AssertExpectations(t)
			mockDocGen.AssertExpectations(t)
		})
	}
}
