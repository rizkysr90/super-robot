package productservice

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store/mocks"
	"rizkysr90-pos/pkg/errorHandler"
)

func TestServiceDeleteProductByID(t *testing.T) {
	tests := []struct {
		name             string
		input            *payload.ReqDeleteProductByID
		mockExpectations func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock)
		expectedError    error
		expectedResult   *payload.ResDeleteProductByID
	}{
		{
			name: "Success - Valid Product ID",
			input: &payload.ReqDeleteProductByID{
				ProductID: "PROD123",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				mockProductStore.On("DeleteByID", mock.Anything, "PROD123").
					Return(nil)
				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
			},
			expectedError:  nil,
			expectedResult: &payload.ResDeleteProductByID{},
		},
		{
			name: "Error - Product Not Found",
			input: &payload.ReqDeleteProductByID{
				ProductID: "NONEXISTENT",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				mockProductStore.On("DeleteByID", mock.Anything, "NONEXISTENT").
					Return(sql.ErrNoRows)
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
			},
			expectedError: errorHandler.NewNotFound(
				errorHandler.WithInfo("product not found"),
			),
			expectedResult: nil,
		},
		{
			name: "Error - Database Error",
			input: &payload.ReqDeleteProductByID{
				ProductID: "PROD123",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				mockProductStore.On("DeleteByID", mock.Anything, "PROD123").
					Return(sql.ErrConnDone)
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
			},
			expectedError: errorHandler.NewInternalServer(
				errorHandler.WithInfo("failed to delete product"),
				errorHandler.WithMessage(sql.ErrConnDone.Error()),
			),
			expectedResult: nil,
		},
		{
			name: "Error - Empty Product ID",
			input: &payload.ReqDeleteProductByID{
				ProductID: "",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				// No mock expectations needed as validation should fail
			},
			expectedError: errorHandler.NewBadRequest(
				errorHandler.WithInfo("validation error"),
				errorHandler.WithMessage("ProductID is required"),
			),
			expectedResult: nil,
		},
		{
			name: "Success - Trims Whitespace",
			input: &payload.ReqDeleteProductByID{
				ProductID: "  PROD123  ",
			},
			mockExpectations: func(mockProductStore *mocks.MockProductStore, sqlMock sqlmock.Sqlmock) {
				mockProductStore.On("DeleteByID", mock.Anything, "PROD123").
					Return(nil)
				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
			},
			expectedError:  nil,
			expectedResult: &payload.ResDeleteProductByID{},
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
			result, err := svc.DeleteProductByID(context.Background(), tt.input)

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
			assert.NoError(t, sqlMock.ExpectationsWereMet())
			mockProductStore.AssertExpectations(t)
		})
	}
}
