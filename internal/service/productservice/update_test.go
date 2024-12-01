package productservice

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/store/mocks"
	"rizkysr90-pos/pkg/errorHandler"
)

func TestServiceUpdateProduct(t *testing.T) {
	// Generate valid UUID for testing
	validCategoryID := uuid.New().String()

	tests := []struct {
		name             string
		input            *payload.ReqUpdateProduct
		mockExpectations func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock)
		expectedError    error
		expectedResult   *payload.ResUpdateProduct
	}{
		{
			name: "Success - Valid Update",
			input: &payload.ReqUpdateProduct{
				ProductID:     "PROD123",
				ProductName:   "Updated Product",
				CategoryID:    validCategoryID,
				Price:         150.0,
				BasePrice:     120.0,
				StockQuantity: 75,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				mockProductStore.On("GetByName", mock.Anything, "UPDATED PRODUCT").
					Return(nil, sql.ErrNoRows)

				mockProductStore.On("Update", mock.Anything, mock.MatchedBy(func(p *store.ProductData) bool {
					return p.ProductID == "PROD123" &&
						p.ProductName == "UPDATED PRODUCT" &&
						p.CategoryID == validCategoryID &&
						p.Price == 150.0 &&
						p.BasePrice == 120.0 &&
						p.StockQuantity == 75
				})).Return(nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
			},
			expectedError:  nil,
			expectedResult: &payload.ResUpdateProduct{},
		},
		{
			name: "Error - Duplicate Product Name",
			input: &payload.ReqUpdateProduct{
				ProductID:     "PROD123",
				ProductName:   "Existing Product",
				CategoryID:    validCategoryID,
				Price:         150.0,
				BasePrice:     120.0,
				StockQuantity: 75,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				mockProductStore.On("GetByName", mock.Anything, "EXISTING PRODUCT").
					Return(&store.ProductData{
						ProductID:   "PROD456", // Different ProductID
						ProductName: "EXISTING PRODUCT",
					}, nil)
			},
			expectedError: errorHandler.NewBadRequest(
				errorHandler.WithInfo("duplicate product name : EXISTING PRODUCT"),
			),
			expectedResult: nil,
		},
		{
			name: "Error - Invalid Category ID Format",
			input: &payload.ReqUpdateProduct{
				ProductID:     "PROD123",
				ProductName:   "Test Product",
				CategoryID:    "invalid-uuid",
				Price:         150.0,
				BasePrice:     120.0,
				StockQuantity: 75,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				// No mocks needed as validation should fail
			},
			expectedError: errorHandler.NewMultipleFieldsValidation([]errorHandler.HttpError{
				{Code: 400, Info: "Validation Error", Message: "CategoryID must be a valid UUID"},
			}),
			expectedResult: nil,
		},
		{
			name: "Error - Empty Required Fields",
			input: &payload.ReqUpdateProduct{
				ProductID:     "",
				ProductName:   "",
				CategoryID:    "",
				Price:         -1,
				BasePrice:     -1,
				StockQuantity: -1,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				// No mocks needed as validation should fail
			},
			expectedError: errorHandler.NewMultipleFieldsValidation([]errorHandler.HttpError{
				{Code: 400, Info: "Validation Error", Message: "ProductID is required"},
				{Code: 400, Info: "Validation Error", Message: "ProductName is required"},
				{Code: 400, Info: "Validation Error", Message: "CategoryID is required"},
				{Code: 400, Info: "Validation Error", Message: "Price must be greater than or equal to 0"},
				{Code: 400, Info: "Validation Error", Message: "BasePrice must be greater than or equal to 0"},
				{Code: 400, Info: "Validation Error", Message: "StockQuantity must be greater than or equal to 0"},
			}),
			expectedResult: nil,
		},
		{
			name: "Success - Same Product ID Different Name",
			input: &payload.ReqUpdateProduct{
				ProductID:     "PROD123",
				ProductName:   "Updated Name",
				CategoryID:    validCategoryID,
				Price:         150.0,
				BasePrice:     120.0,
				StockQuantity: 75,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				// Product exists but has same ID (allowed)
				mockProductStore.On("GetByName", mock.Anything, "UPDATED NAME").
					Return(&store.ProductData{
						ProductID:   "PROD123", // Same ProductID
						ProductName: "UPDATED NAME",
					}, nil)

				mockProductStore.On("Update", mock.Anything, mock.AnythingOfType("*store.ProductData")).
					Return(nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
			},
			expectedError:  nil,
			expectedResult: &payload.ResUpdateProduct{},
		},
		{
			name: "Error - Database Error on Update",
			input: &payload.ReqUpdateProduct{
				ProductID:     "PROD123",
				ProductName:   "Test Product",
				CategoryID:    validCategoryID,
				Price:         150.0,
				BasePrice:     120.0,
				StockQuantity: 75,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				mockProductStore.On("GetByName", mock.Anything, "TEST PRODUCT").
					Return(nil, sql.ErrNoRows)

				mockProductStore.On("Update", mock.Anything, mock.AnythingOfType("*store.ProductData")).
					Return(sql.ErrConnDone)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
			},
			expectedError:  fmt.Errorf("sqldb: WithinTx failed before commit: %w", sql.ErrConnDone),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			db, sqlMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mockProductStore := new(mocks.MockProductStore)
			tt.mockExpectations(mockProductStore, sqlMock)

			svc := &Service{
				db:           db,
				productStore: mockProductStore,
			}

			// Execute
			result, err := svc.UpdateProduct(context.Background(), tt.input)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			// Verify mock expectations
			assert.NoError(t, sqlMock.ExpectationsWereMet())
			mockProductStore.AssertExpectations(t)
		})
	}
}
