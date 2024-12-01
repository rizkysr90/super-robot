package productservice

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/store/mocks"
	"rizkysr90-pos/pkg/errorHandler"
)

func TestServiceCreateProduct(t *testing.T) {
	tests := []struct {
		name             string
		input            *payload.ReqCreateProduct
		mockExpectations func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock)
		expectedError    error
		expectedProduct  *payload.ResCreateProduct
	}{
		{
			name: "Success - Valid Input",
			input: &payload.ReqCreateProduct{
				ProductName:   "Test Product",
				CategoryID:    "CAT123",
				Price:         100.0,
				BasePrice:     80.0,
				StockQuantity: 50,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				mockProductStore.On("GetByName", mock.Anything, "TEST PRODUCT").
					Return(nil, sql.ErrNoRows)

				mockProductStore.On("Insert", mock.Anything, mock.MatchedBy(func(p *store.ProductData) bool {
					return p.ProductName == "TEST PRODUCT" &&
						p.CategoryID == "CAT123" &&
						p.Price == 100.0 &&
						p.BasePrice == 80.0 &&
						p.StockQuantity == 50
				})).Return(nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
			},
			expectedError:   nil,
			expectedProduct: &payload.ResCreateProduct{},
		},
		{
			name: "Error - Duplicate Product Name",
			input: &payload.ReqCreateProduct{
				ProductName:   "Existing Product",
				CategoryID:    "CAT123",
				Price:         100.0,
				BasePrice:     80.0,
				StockQuantity: 50,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				existingProduct := &store.ProductData{
					ProductName: "EXISTING PRODUCT",
					ProductID:   "PROD123",
				}
				mockProductStore.On("GetByName", mock.Anything, "EXISTING PRODUCT").
					Return(existingProduct, nil)
			},
			expectedError:   errorHandler.NewBadRequest(errorHandler.WithInfo("duplicate product name : EXISTING PRODUCT")),
			expectedProduct: nil,
		},
		{
			name: "Error - Invalid Input (Missing Required Fields)",
			input: &payload.ReqCreateProduct{
				ProductName:   "",
				CategoryID:    "",
				Price:         -1,
				BasePrice:     -1,
				StockQuantity: -1,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				// No mocks needed as validation should fail before DB calls
			},
			expectedError: errorHandler.NewMultipleFieldsValidation([]errorHandler.HttpError{
				{Code: 400, Info: "Validation Error", Message: "ProductName is required"},
				{Code: 400, Info: "Validation Error", Message: "CategoryID is required"},
				{Code: 400, Info: "Validation Error", Message: "Price must be greater than or equal to 0"},
				{Code: 400, Info: "Validation Error", Message: "BasePrice must be greater than or equal to 0"},
				{Code: 400, Info: "Validation Error", Message: "StockQuantity must be greater than or equal to 0"},
			}),
			expectedProduct: nil,
		},
		{
			name: "Error - Database Error on GetByName",
			input: &payload.ReqCreateProduct{
				ProductName:   "Test Product",
				CategoryID:    "CAT123",
				Price:         100.0,
				BasePrice:     80.0,
				StockQuantity: 50,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				mockProductStore.On("GetByName", mock.Anything, "TEST PRODUCT").
					Return(nil, sql.ErrConnDone)
			},
			expectedError:   sql.ErrConnDone,
			expectedProduct: nil,
		},
		{
			name: "Error - Database Error on Insert",
			input: &payload.ReqCreateProduct{
				ProductName:   "Test Product",
				CategoryID:    "CAT123",
				Price:         100.0,
				BasePrice:     80.0,
				StockQuantity: 50,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				mockProductStore.On("GetByName", mock.Anything, "TEST PRODUCT").
					Return(nil, sql.ErrNoRows)

				mockProductStore.On("Insert", mock.Anything, mock.AnythingOfType("*store.ProductData")).
					Return(sql.ErrConnDone)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
			},
			expectedError:   fmt.Errorf("sqldb: WithinTx failed before commit: %w", sql.ErrConnDone),
			expectedProduct: nil,
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
			result, err := svc.CreateProduct(context.Background(), tt.input)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedProduct, result)
			}

			// Verify all mock expectations were met
			assert.NoError(t, sqlMock.ExpectationsWereMet())
			mockProductStore.AssertExpectations(t)
		})
	}
}
