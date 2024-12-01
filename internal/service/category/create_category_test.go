package category

import (
	"context"
	"errors"
	"testing"

	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/store/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCategoryService(t *testing.T) {
	tests := []struct {
		name             string
		input            *payload.ReqCreateCategory
		mockExpectations func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock)
		expectedError    func(err error) bool
		expectedCategory *payload.ResCreateCategory
	}{
		{
			name: "Valid Input",
			input: &payload.ReqCreateCategory{
				CategoryName: "New Category",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByName", context.Background(), "NEW CATEGORY").Return(nil, nil)
				mockCategoryStore.On("Create", mock.Anything, mock.MatchedBy(func(categoryData *store.CategoryData) bool {
					return categoryData.CategoryName == "NEW CATEGORY" // Adjust to your actual logic
				})).Return(nil)
				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
			},
			expectedError:    nil,
			expectedCategory: &payload.ResCreateCategory{},
		},
		{
			name: "Empty Category Name",
			input: &payload.ReqCreateCategory{
				CategoryName: "",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				// it empty because if category name it should run any assertion
			},
			expectedError: func(err error) bool {
				// Check if the error message contains the expected prefix
				return err != nil && err.Error() == "category name is required"
			},
			expectedCategory: nil,
		},
		{
			name: "Category Name Too Long",
			input: &payload.ReqCreateCategory{
				CategoryName: "A very long category name that exceeds the maximum allowed length for a category name  allowed length for a category name",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				// it empty because if category name it should run any assertion
			},
			expectedError: func(err error) bool {
				// Check if the error message contains the expected prefix
				return err != nil && err.Error() == "max category name is 100 characters"
			},
			expectedCategory: nil,
		},
		{
			name: "Duplicate Category Name",
			input: &payload.ReqCreateCategory{
				CategoryName: "Existing Category",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByName", context.Background(), "EXISTING CATEGORY").Return(&store.CategoryData{}, nil)
			},
			expectedError: func(err error) bool {
				// Check if the error message contains the expected prefix
				return err != nil && err.Error() == "duplicate category name"
			},
			expectedCategory: nil,
		},
		{
			name: "Database Error During Creation",
			input: &payload.ReqCreateCategory{
				CategoryName: "Valid Category",
			},
			mockExpectations: func(mockCategoryStore *mocks.MockCategoryStore, sqlMock sqlmock.Sqlmock) {
				mockCategoryStore.On("FindByName", context.Background(), "VALID CATEGORY").Return(nil, nil)
				mockCategoryStore.On("Create", mock.Anything, mock.MatchedBy(func(categoryData *store.CategoryData) bool {
					return categoryData.CategoryName == "VALID CATEGORY"
				})).Return(errors.New("db error"))
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
			},
			expectedError: func(err error) bool {
				// Check if the error message contains the expected prefix
				return err != nil && err.Error() == "sqldb: WithinTx failed before commit: db error"
			},
			expectedCategory: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, mockCategoryStore, sqlMock, _ := initTestService(t)

			tt.mockExpectations(mockCategoryStore, sqlMock)

			res, err := svc.Create(context.Background(), tt.input)

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
