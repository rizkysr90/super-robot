package category

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/store/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCategories(t *testing.T) {
	fixedTime := time.Now().UTC()
	tests := []struct {
		name             string
		input            *payload.ReqGetAllCategory
		mockExpectations func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock)
		expectedError    func(err error) bool
		expectedResponse *payload.ResGetAllCategory
	}{
		{
			name: "Valid Input with Categories",
			input: &payload.ReqGetAllCategory{
				PageNumber: 1,
				PageSize:   10,
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindAllPagination", mock.Anything, &store.Pagination{
					PageSize:   10,
					PageNumber: 1,
				}).Return([]store.CategoryData{
					{
						ID:           "1",
						CategoryName: "Category 1",
						CreatedAt:    fixedTime,
						UpdatedAt:    fixedTime,
						DeletedAt:    sql.NullTime{},
						Pagination: &store.Pagination{
							TotalElements: 2,
						},
					},
					{
						ID:           "2",
						CategoryName: "Category 2",
						CreatedAt:    fixedTime,
						UpdatedAt:    fixedTime,
						DeletedAt:    sql.NullTime{},
						Pagination: &store.Pagination{
							TotalElements: 2,
						},
					},
				}, nil)
			},
			expectedError: nil,
			expectedResponse: &payload.ResGetAllCategory{
				Data: []payload.CategoryData{
					{
						ID:           "1",
						CategoryName: "Category 1",
						CreatedAt:    fixedTime,
						UpdatedAt:    fixedTime,
						DeletedAt:    time.Time{},
					},
					{
						ID:           "2",
						CategoryName: "Category 2",
						CreatedAt:    fixedTime,
						UpdatedAt:    fixedTime,
						DeletedAt:    time.Time{},
					},
				},
				Metadata: payload.Pagination{
					PageSize:      10,
					PageNumber:    1,
					TotalPages:    1,
					TotalElements: 2,
				},
			},
		},
		{
			name: "No Categories Found",
			input: &payload.ReqGetAllCategory{
				PageNumber: 1,
				PageSize:   10,
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindAllPagination", mock.Anything, &store.Pagination{
					PageSize:   10,
					PageNumber: 1,
				}).Return(nil, nil)
			},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "not found"
			},
			expectedResponse: nil,
		},
		{
			name: "Database Error During Retrieval",
			input: &payload.ReqGetAllCategory{
				PageNumber: 1,
				PageSize:   10,
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindAllPagination", mock.Anything, &store.Pagination{
					PageSize:   10,
					PageNumber: 1,
				}).Return(nil, errors.New("db error"))
			},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "db error"
			},
			expectedResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the service and mocks
			svc, mockCategoryStore, sqlMock, _ := initTestService(t)

			// Set up mock expectations
			tt.mockExpectations(mockCategoryStore, sqlMock)

			// Call the service method
			resp, err := svc.GetCategories(context.Background(), tt.input)

			// Assert the results
			if tt.expectedError != nil {
				assert.True(t, tt.expectedError(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, resp)
			}

			// Assert SQL expectations
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}
