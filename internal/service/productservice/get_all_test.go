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

func TestServiceGetAllProducts(t *testing.T) {
	// Define fixed time for consistent testing
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name             string
		input            *payload.ReqGetAllProducts
		mockExpectations func(mockProductStore *mocks.MockProductStore)
		expectedResult   *payload.ResGetAllProducts
		expectedError    error
	}{
		{
			name: "Success - Default Pagination",
			input: &payload.ReqGetAllProducts{
				PageSize:   0, // Should default to 10
				PageNumber: 0, // Should default to 1
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore) {
				mockProductStore.On("GetAll", mock.Anything, &store.FilterProduct{
					Limit:      10,
					Offset:     0,
					CategoryID: "",
				}).Return([]store.ProductData{
					{
						ProductID:     "PROD1",
						ProductName:   "Test Product 1",
						Price:         100.0,
						BasePrice:     80.0,
						StockQuantity: 50,
						Category: &store.CategoryData{
							ID:           "CAT1",
							CategoryName: "Test Category",
						},
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					},
				}, 1, nil)
			},
			expectedResult: &payload.ResGetAllProducts{
				Data: []payload.ProductData{
					{
						ProductID:     "PROD1",
						ProductName:   "Test Product 1",
						Price:         100.0,
						BasePrice:     80.0,
						StockQuantity: 50,
						CategoryID:    "CAT1",
						CategoryName:  "Test Category",
						CreatedAt:     fixedTime,
						UpdatedAt:     fixedTime,
					},
				},
				Metadata: payload.Pagination{
					PageSize:      10,
					PageNumber:    1,
					TotalPages:    1,
					TotalElements: 1,
				},
			},
			expectedError: nil,
		},
		{
			name: "Success - Custom Pagination and Category Filter",
			input: &payload.ReqGetAllProducts{
				PageSize:   5,
				PageNumber: 2,
				CategoryID: "CAT1",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore) {
				mockProductStore.On("GetAll", mock.Anything, &store.FilterProduct{
					Limit:      5,
					Offset:     5, // (PageNumber-1) * PageSize
					CategoryID: "CAT1",
				}).Return([]store.ProductData{
					{
						ProductID:     "PROD2",
						ProductName:   "Test Product 2",
						Price:         200.0,
						BasePrice:     160.0,
						StockQuantity: 30,
						Category: &store.CategoryData{
							ID:           "CAT1",
							CategoryName: "Test Category",
						},
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					},
				}, 8, nil) // Total 8 items
			},
			expectedResult: &payload.ResGetAllProducts{
				Data: []payload.ProductData{
					{
						ProductID:     "PROD2",
						ProductName:   "Test Product 2",
						Price:         200.0,
						BasePrice:     160.0,
						StockQuantity: 30,
						CategoryID:    "CAT1",
						CategoryName:  "Test Category",
						CreatedAt:     fixedTime,
						UpdatedAt:     fixedTime,
					},
				},
				Metadata: payload.Pagination{
					PageSize:      5,
					PageNumber:    2,
					TotalPages:    2, // Ceil(8/5)
					TotalElements: 8,
				},
			},
			expectedError: nil,
		},
		{
			name: "Success - Empty Result",
			input: &payload.ReqGetAllProducts{
				PageSize:   10,
				PageNumber: 1,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore) {
				mockProductStore.On("GetAll", mock.Anything, &store.FilterProduct{
					Limit:      10,
					Offset:     0,
					CategoryID: "",
				}).Return([]store.ProductData{}, 0, nil)
			},
			expectedResult: &payload.ResGetAllProducts{
				Data: []payload.ProductData{},
				Metadata: payload.Pagination{
					PageSize:      10,
					PageNumber:    1,
					TotalPages:    0,
					TotalElements: 0,
				},
			},
			expectedError: nil,
		},
		{
			name: "Error - Database Error",
			input: &payload.ReqGetAllProducts{
				PageSize:   10,
				PageNumber: 1,
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore) {
				mockProductStore.On("GetAll", mock.Anything, &store.FilterProduct{
					Limit:      10,
					Offset:     0,
					CategoryID: "",
				}).Return(nil, 0, sql.ErrConnDone)
			},
			expectedResult: nil,
			expectedError: errorHandler.NewInternalServer(
				errorHandler.WithInfo("failed to get products"),
				errorHandler.WithMessage(sql.ErrConnDone.Error()),
			),
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
			result, err := svc.GetAllProducts(context.Background(), tt.input)

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
