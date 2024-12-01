package category

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/store/mocks" // Adjust the import path as necessary

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCategory(t *testing.T) {
	const existingID = "existing-id"
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
				ID: existingID,
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByID", mock.Anything, existingID).Return(&store.CategoryData{
					ID: existingID,
				}, nil)
				mockCategoryStore.On("SoftDelete", mock.Anything, existingID).Return(nil)
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
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				// it's empty because we dont need mocking
			},
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
				ID: existingID,
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByID", mock.Anything, existingID).Return(&store.CategoryData{
					ID: existingID,
				}, nil)
				mockCategoryStore.On("SoftDelete", mock.Anything, existingID).Return(errors.New("db error"))
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
