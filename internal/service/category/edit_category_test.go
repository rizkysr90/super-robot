package category

import (
	"context"
	"errors"
	"strings"
	"testing"

	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/store/mocks" // Import your store mocks

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceEditCategory(t *testing.T) {
	fixedUUID := uuid.NewString()
	tests := []struct {
		name             string
		input            *payload.ReqUpdateCategory
		mockExpectations func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock)
		expectedError    func(err error) bool
		expectedCategory *payload.ResUpdateCategory
	}{
		{
			name: "Valid Input",
			input: &payload.ReqUpdateCategory{
				ID:           fixedUUID, // Use a fixed value if you want consistency
				CategoryName: "Updated Category",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				// Mock the FindByName method to return a category with the same ID
				mockCategoryStore.On("FindByName", mock.Anything, "UPDATED CATEGORY").Return(&store.CategoryData{
					ID: fixedUUID, // This should be a fixed value to match the input ID
				}, nil)

				// Mock the Update method to check that the ID matches and return no error
				mockCategoryStore.On("Update", mock.Anything, mock.MatchedBy(func(categoryData *store.CategoryData) bool {
					return categoryData.ID == fixedUUID && categoryData.CategoryName == "UPDATED CATEGORY"
				})).Return(nil)

				// Set up SQL expectations for transaction
				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
			},
			expectedError:    nil,
			expectedCategory: &payload.ResUpdateCategory{},
		},
		{
			name: "Validation Error - Missing ID",
			input: &payload.ReqUpdateCategory{
				CategoryName: "Valid Category",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				// it's empty because we dont need mocking
			},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "id is required"
			},
			expectedCategory: nil,
		},
		{
			name: "Validation Error - Missing Category Name",
			input: &payload.ReqUpdateCategory{
				ID: "valid-id",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				// it's empty because we dont need mocking
			},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "category name is required"
			},
			expectedCategory: nil,
		},
		{
			name: "Validation Error - Category Name Too Long",
			input: &payload.ReqUpdateCategory{
				ID:           "valid-id",
				CategoryName: strings.Repeat("A", 201), // Exceeding the limit
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				// it's empty because we dont need mocking
			},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "too long, max 200 character"
			},
			expectedCategory: nil,
		},
		{
			name: "Duplicate Category Name",
			input: &payload.ReqUpdateCategory{
				ID:           uuid.NewString(),
				CategoryName: "Existing Category",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByName", mock.Anything, "EXISTING CATEGORY").Return(&store.CategoryData{
					ID: "existing-id",
				}, nil)
			},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "duplicate category name"
			},
			expectedCategory: nil,
		},
		{
			name: "Database Error During Update",
			input: &payload.ReqUpdateCategory{
				ID:           fixedUUID,
				CategoryName: "Valid Category",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByName", mock.Anything, "VALID CATEGORY").Return(&store.CategoryData{
					ID: fixedUUID,
				}, nil)
				mockCategoryStore.On("Update", mock.Anything, mock.MatchedBy(func(categoryData *store.CategoryData) bool {
					return categoryData.ID == fixedUUID && categoryData.CategoryName == "VALID CATEGORY"
				})).Return(errors.New("db error"))

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
			},
			expectedError: func(err error) bool {
				return err != nil && err.Error() == "sqldb: WithinTx failed before commit: db error"
			},
			expectedCategory: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, mockCategoryStore, sqlMock, _ := initTestService(t)

			tt.mockExpectations(mockCategoryStore, sqlMock)

			res, err := svc.EditCategory(context.Background(), tt.input)

			if tt.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCategory, res)
			} else {
				assert.True(t, tt.expectedError(err))
				assert.Nil(t, res)
			}

			mockCategoryStore.AssertExpectations(t)
			if err := sqlMock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
