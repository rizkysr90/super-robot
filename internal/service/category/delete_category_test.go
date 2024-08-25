package category

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/internal/store"
	"auth-service-rizkysr90-pos/internal/store/mocks" // Adjust the import path as necessary

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCategory(t *testing.T) {
	tests := []struct {
		name             string
		input            *payload.ReqDeleteCategory
		mockExpectations func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock)
		expectedError    func(err error) bool
		expectedResponse *payload.ResDeleteCategory
	}{
		{
			name: "Valid Input",
			input: &payload.ReqDeleteCategory{
				ID: "existing-id",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByID", mock.Anything, "existing-id").Return(&store.CategoryData{
					ID: "existing-id",
				}, nil)
				mockCategoryStore.On("SoftDelete", mock.Anything, "existing-id").Return(nil)
				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
			},
			expectedError:    nil,
			expectedResponse: &payload.ResDeleteCategory{},
		},
		{
			name: "Validation Error - Missing ID",
			input: &payload.ReqDeleteCategory{
				ID: "",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "required id data"
			},
			expectedResponse: nil,
		},
		{
			name: "Category Not Found",
			input: &payload.ReqDeleteCategory{
				ID: "non-existing-id",
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
			name: "Database Error During Deletion",
			input: &payload.ReqDeleteCategory{
				ID: "existing-id",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByID", mock.Anything, "existing-id").Return(&store.CategoryData{
					ID: "existing-id",
				}, nil)
				mockCategoryStore.On("SoftDelete", mock.Anything, "existing-id").Return(errors.New("db error"))
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
			},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "sqldb: WithinTx failed before commit: db error"
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
			resp, err := svc.DeleteCategory(context.Background(), tt.input)

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
