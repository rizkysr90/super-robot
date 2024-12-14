package auth

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/store/mocks"
	"rizkysr90-pos/pkg/errorHandler"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterOwner(t *testing.T) {
	tests := []struct {
		name             string
		input            *RequestRegisterOwner
		mockExpectations func(*mocks.MockStateStore, sqlmock.Sqlmock)
		expectedError    error
		expectedStateID  string
	}{
		{
			name: "Valid Input",
			input: &RequestRegisterOwner{
				TenantName: "Test Tenant",
			},
			mockExpectations: func(mockStateStore *mocks.MockStateStore, sqlMock sqlmock.Sqlmock) {
				mockStateStore.On("Insert", mock.Anything, mock.MatchedBy(func(stateData *store.StateData) bool {
					return stateData.TenantName.String == "Test Tenant" &&
						stateData.TenantName.Valid &&
						len(stateData.ID) == 24
				})).Return(nil)
				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
			},
			expectedError:   nil,
			expectedStateID: "", // Will be random, we'll just check the length
		},
		{
			name: "Empty Tenant Name",
			input: &RequestRegisterOwner{
				TenantName: "",
			},
			mockExpectations: func(mockStateStore *mocks.MockStateStore, sqlMock sqlmock.Sqlmock) {
				// No expectations needed as validation should fail first
			},
			expectedError: errorHandler.NewBadRequest(errorHandler.WithInfo(
				"failed to validate request, got : Bad Request - Error Multiple Fields Validaition",
			)),
			expectedStateID: "",
		},
		{
			name: "Database Insert Error",
			input: &RequestRegisterOwner{
				TenantName: "Test Tenant",
			},
			mockExpectations: func(mockStateStore *mocks.MockStateStore, sqlMock sqlmock.Sqlmock) {
				mockStateStore.On("Insert", mock.Anything, mock.Anything).
					Return(sql.ErrNoRows)
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
			},
			expectedError: errorHandler.NewInternalServer(errorHandler.WithInfo(
				"failed to insert state data, got : sqldb: WithinTx failed before commit: sql: no rows in result set",
			)),
			expectedStateID: "",
		},
		{
			name: "Transaction Begin Error",
			input: &RequestRegisterOwner{
				TenantName: "Test Tenant",
			},
			mockExpectations: func(mockStateStore *mocks.MockStateStore, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
			expectedError: errorHandler.NewInternalServer(errorHandler.WithInfo(
				"failed to insert state data, got : sqldb: WithinTx begin SQL transaction failed: sql: connection is already closed",
			)),
			expectedStateID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockStateStore := &mocks.MockStateStore{}
			db, sqlMock, err := sqlmock.New()
			if err != nil {
				panic(err)
			}
			defer db.Close()

			// Setup mock expectations
			tt.mockExpectations(mockStateStore, sqlMock)

			// Create auth instance
			auth := &Auth{
				stateStore: mockStateStore,
				db:         db,
			}

			// Execute test
			stateID, err := auth.RegisterOwner(context.Background(), tt.input)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, stateID)
			} else {
				assert.NoError(t, err)
				assert.Len(t, stateID, 24) // base64 encoded 16 bytes
			}

			// Verify mock expectations
			mockStateStore.AssertExpectations(t)
			assert.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}
