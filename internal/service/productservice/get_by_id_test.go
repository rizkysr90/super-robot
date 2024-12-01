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

func TestServiceGetProductByID(t *testing.T) {
	// Define fixed time for consistent testing
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		input            *payload.ReqGetProductByID
		mockExpectations func(mockProductStore *mocks.MockProductStore)
		expectedResult   *payload.ResGetProductByID
		expectedError    error
	}{
		{
			name: "Success - Valid Product ID",
			input: &payload.ReqGetProductByID{
				ProductID: "PROD123",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore) {
				mockProductStore.On("GetByID", mock.Anything, "PROD123").
					Return(&store.ProductData{
						ProductID:     "PROD123",
						ProductName:   "Test Product",
						Price:         100.0,
						BasePrice:     80.0,
						StockQuantity: 50,
						CategoryID:    "CAT1",
						Category: &store.CategoryData{
							ID:           "CAT1",
							CategoryName: "Test Category",
						},
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					}, nil)
			},
			expectedResult: &payload.ResGetProductByID{
				Data: payload.ProductData{
					ProductID:     "PROD123",
					ProductName:   "Test Product",
					Price:         100.0,
					BasePrice:     80.0,
					StockQuantity: 50,
					CategoryID:    "CAT1",
					CategoryName:  "Test Category",
					CreatedAt:     fixedTime,
					UpdatedAt:     fixedTime,
				},
			},
			expectedError: nil,
		},
		{
			name: "Success - Product ID with Whitespace",
			input: &payload.ReqGetProductByID{
				ProductID: "  PROD123  ",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore) {
				mockProductStore.On("GetByID", mock.Anything, "PROD123").
					Return(&store.ProductData{
						ProductID:     "PROD123",
						ProductName:   "Test Product",
						Price:         100.0,
						BasePrice:     80.0,
						StockQuantity: 50,
						CategoryID:    "CAT1",
						Category: &store.CategoryData{
							ID:           "CAT1",
							CategoryName: "Test Category",
						},
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					}, nil)
			},
			expectedResult: &payload.ResGetProductByID{
				Data: payload.ProductData{
					ProductID:     "PROD123",
					ProductName:   "Test Product",
					Price:         100.0,
					BasePrice:     80.0,
					StockQuantity: 50,
					CategoryID:    "CAT1",
					CategoryName:  "Test Category",
					CreatedAt:     fixedTime,
					UpdatedAt:     fixedTime,
				},
			},
			expectedError: nil,
		},
		{
			name: "Error - Product Not Found",
			input: &payload.ReqGetProductByID{
				ProductID: "NONEXISTENT",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore) {
				mockProductStore.On("GetByID", mock.Anything, "NONEXISTENT").
					Return(nil, sql.ErrNoRows)
			},
			expectedResult: nil,
			expectedError: errorHandler.NewBadRequest(
				errorHandler.WithInfo("product id not found"),
			),
		},
		{
			name: "Error - Database Error",
			input: &payload.ReqGetProductByID{
				ProductID: "PROD123",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore) {
				mockProductStore.On("GetByID", mock.Anything, "PROD123").
					Return(nil, sql.ErrConnDone)
			},
			expectedResult: nil,
			expectedError:  sql.ErrConnDone,
		},
		{
			name: "Error - Empty Product ID",
			input: &payload.ReqGetProductByID{
				ProductID: "",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore) {
				// No mock expectations needed as validation should fail
			},
			expectedResult: nil,
			expectedError: errorHandler.NewBadRequest(
				errorHandler.WithInfo("validation error"),
				errorHandler.WithMessage("Validation failed: ProductID is required"),
			),
		},
		{
			name: "Success - Soft Deleted Product",
			input: &payload.ReqGetProductByID{
				ProductID: "PROD123",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore) {
				deletedTime := fixedTime.Add(time.Hour * 24)
				mockProductStore.On("GetByID", mock.Anything, "PROD123").
					Return(&store.ProductData{
						ProductID:     "PROD123",
						ProductName:   "Test Product",
						Price:         100.0,
						BasePrice:     80.0,
						StockQuantity: 50,
						CategoryID:    "CAT1",
						Category: &store.CategoryData{
							ID:           "CAT1",
							CategoryName: "Test Category",
						},
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: deletedTime,
					}, nil)
			},
			expectedResult: &payload.ResGetProductByID{
				Data: payload.ProductData{
					ProductID:     "PROD123",
					ProductName:   "Test Product",
					Price:         100.0,
					BasePrice:     80.0,
					StockQuantity: 50,
					CategoryID:    "CAT1",
					CategoryName:  "Test Category",
					CreatedAt:     fixedTime,
					UpdatedAt:     fixedTime,
					DeletedAt:     fixedTime.Add(time.Hour * 24),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockProductStore := new(mocks.MockProductStore)
			tt.mockExpectations(mockProductStore)

			svc := &Service{
				productStore: mockProductStore,
			}

			// Execute
			result, err := svc.GetProductByID(context.Background(), tt.input)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			// Verify mock expectations
			mockProductStore.AssertExpectations(t)
		})
	}
}
