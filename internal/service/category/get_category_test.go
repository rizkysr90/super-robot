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

func TestGetCategoryByID(t *testing.T) {
	fixedTime := time.Now().UTC()
	tests := []struct {
		name             string
		input            *payload.ReqGetCategoryByID
		mockExpectations func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock)
		expectedError    func(err error) bool
		expectedResponse *payload.ResGetCategoryByID
	}{
		{
			name: "Valid Input",
			input: &payload.ReqGetCategoryByID{
				CategoryID: "existing-id",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByID", mock.Anything, "existing-id").Return(&store.CategoryData{
					ID:           "existing-id",
					CategoryName: "Category Name",
					CreatedAt:    fixedTime,
					UpdatedAt:    fixedTime,
					DeletedAt:    sql.NullTime{},
				}, nil)
			},
			expectedError: nil,
			expectedResponse: &payload.ResGetCategoryByID{
				CategoryData: &payload.CategoryData{
					ID:           "existing-id",
					CategoryName: "Category Name",
					CreatedAt:    fixedTime,
					UpdatedAt:    fixedTime,
					DeletedAt:    time.Time{},
				},
			},
		},
		{
			name: "Validation Error - Missing ID",
			input: &payload.ReqGetCategoryByID{
				CategoryID: "",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "category_id is required"
			},
			expectedResponse: nil,
		},
		{
			name: "Category Not Found",
			input: &payload.ReqGetCategoryByID{
				CategoryID: "non-existing-id",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByID", mock.Anything, "non-existing-id").Return(nil, sql.ErrNoRows)
			},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "category not found"
			},
			expectedResponse: nil,
		},
		{
			name: "Database Error During Retrieval",
			input: &payload.ReqGetCategoryByID{
				CategoryID: "existing-id",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByID", mock.Anything, "existing-id").Return(nil, errors.New("db error"))
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
			resp, err := svc.GetCategoryByID(context.Background(), tt.input)

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
