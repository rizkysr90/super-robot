package branches

import (
	"context"
	"database/sql"
	"rizkysr90-pos/internal/commonvalidator"
	"rizkysr90-pos/internal/config"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/store/mocks"
	"rizkysr90-pos/internal/utility"
	"rizkysr90-pos/pkg/errorHandler"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		request       *RequestCreate
		setupMocks    func(*mocks.MockUser, *mocks.MockTenant, *mocks.MockBranchStore, sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "success_create_branch",
			request: &RequestCreate{
				TenantID:   "tenant-123",
				BranchName: "BRANCH-1",
				Address:    "TEST ADDRESS",
				CreatedBy:  "user-123",
			},
			setupMocks: func(us *mocks.MockUser, ts *mocks.MockTenant, bs *mocks.MockBranchStore, sqlMock sqlmock.Sqlmock) {
				// User exists and is owner
				us.On("FindOne", mock.Anything, &store.UserQueryFilter{ID: "user-123"}).
					Return(&store.UserData{ID: "user-123"}, nil)

				// Tenant exists and user is owner
				ts.On("FindOne", mock.Anything, &store.TenantFilter{ID: "tenant-123"}).
					Return(&store.TenantData{
						ID:      "tenant-123",
						OwnerID: sql.NullString{String: "user-123", Valid: true},
					}, nil)

				// Branch doesn't exist yet
				bs.On("FindOne", mock.Anything, &store.BranchesFilter{
					Name:     "BRANCH-1",
					TenantID: "tenant-123",
				}).Return(nil, sql.ErrNoRows)

				// Expect transaction begin
				sqlMock.ExpectBegin()

				// Expect branch insert
				bs.On("Insert", mock.Anything, mock.MatchedBy(func(data *store.BranchesData) bool {
					return data.TenantID == "tenant-123" &&
						data.Name == "BRANCH-1" &&
						data.Address == "TEST ADDRESS" &&
						data.CreatedByOwner.String == "user-123"
				})).Return(nil)
				bs.On("TotalBranches", mock.Anything, "tenant-123").Return(1, nil)
				// Expect transaction commit
				sqlMock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "validation_error_missing_fields",
			request: &RequestCreate{
				TenantID:   "",
				BranchName: "",
				Address:    "",
				CreatedBy:  "",
			},
			setupMocks: func(us *mocks.MockUser, ts *mocks.MockTenant, bs *mocks.MockBranchStore, sqlMock sqlmock.Sqlmock) {
				// No mocks needed as validation should fail first
			},
			expectedError: errorHandler.NewMultipleFieldsValidation([]errorHandler.HttpError{
				*utility.ConstructErrorRequired("tenant_id"),
				*utility.ConstructErrorRequired("branch_name"),
				*utility.ConstructErrorRequired("address"),
				*utility.ConstructErrorRequired("created_by"),
			}),
		},
		{
			name: "error_branch_already_exists",
			request: &RequestCreate{
				TenantID:   "tenant-123",
				BranchName: "EXISTING-BRANCH",
				Address:    "TEST ADDRESS",
				CreatedBy:  "user-123",
			},
			setupMocks: func(us *mocks.MockUser, ts *mocks.MockTenant, bs *mocks.MockBranchStore, sqlMock sqlmock.Sqlmock) {
				// User exists and is owner
				us.On("FindOne", mock.Anything, &store.UserQueryFilter{ID: "user-123"}).
					Return(&store.UserData{ID: "user-123"}, nil)

				// Tenant exists and user is owner
				ts.On("FindOne", mock.Anything, &store.TenantFilter{ID: "tenant-123"}).
					Return(&store.TenantData{
						ID:      "tenant-123",
						OwnerID: sql.NullString{String: "user-123", Valid: true},
					}, nil)

				// Branch already exists
				bs.On("FindOne", mock.Anything, &store.BranchesFilter{
					Name:     "EXISTING-BRANCH",
					TenantID: "tenant-123",
				}).Return(&store.BranchesData{}, nil)
			},
			expectedError: errorHandler.NewBadRequest(errorHandler.WithInfo("branch is exist")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock DB and stores
			db, sqlMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			userStore := new(mocks.MockUser)
			tenantStore := new(mocks.MockTenant)
			branchStore := new(mocks.MockBranchStore)

			// Setup mocks
			tt.setupMocks(userStore, tenantStore, branchStore, sqlMock)

			// Create service
			svc := &Service{
				cfg:         &config.Config{MaxBranchesPerTenant: 2},
				db:          db,
				userStore:   userStore,
				tenantStore: tenantStore,
				branchStore: branchStore,
			}

			// Execute test
			err = svc.Create(context.Background(), tt.request)

			// Verify expectations
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify all mocks
			userStore.AssertExpectations(t)
			tenantStore.AssertExpectations(t)
			branchStore.AssertExpectations(t)
			assert.NoError(t, sqlMock.ExpectationsWereMet())
		})
	}
}

func TestReqCreate_Sanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    *reqCreate
		expected *reqCreate
	}{
		{
			name: "sanitize_with_spaces_and_lowercase",
			input: &reqCreate{
				RequestCreate: &RequestCreate{
					TenantID:   "  tenant-123  ",
					BranchName: "  branch-1  ",
					Address:    "  test address  ",
					CreatedBy:  "  user-123  ",
				},
			},
			expected: &reqCreate{
				RequestCreate: &RequestCreate{
					TenantID:   "tenant-123",
					BranchName: "BRANCH-1",
					Address:    "TEST ADDRESS",
					CreatedBy:  "user-123",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.sanitize()
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

func TestIsAllowedUser(t *testing.T) {
	tests := []struct {
		name              string
		userID            string
		setupMock         func(*mocks.MockUser)
		expectedIsAllowed bool
		expectedUser      *store.UserData
		expectedError     error
	}{
		{
			name:   "user_exists",
			userID: "user-123",
			setupMock: func(us *mocks.MockUser) {
				us.On("FindOne", mock.Anything, &store.UserQueryFilter{ID: "user-123"}).
					Return(&store.UserData{ID: "user-123"}, nil)
			},
			expectedIsAllowed: true,
			expectedUser:      &store.UserData{ID: "user-123"},
			expectedError:     nil,
		},
		{
			name:   "user_not_found",
			userID: "non-existent",
			setupMock: func(us *mocks.MockUser) {
				us.On("FindOne", mock.Anything, &store.UserQueryFilter{ID: "non-existent"}).
					Return(nil, sql.ErrNoRows)
			},
			expectedIsAllowed: false,
			expectedUser:      nil,
			expectedError:     errorHandler.NewUnauthorized(errorHandler.WithInfo("not allowed")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStore := new(mocks.MockUser)
			tt.setupMock(userStore)

			isAllowed, userData, err := commonvalidator.IsAllowedUser(context.Background(), userStore, tt.userID)

			assert.Equal(t, tt.expectedIsAllowed, isAllowed)
			assert.Equal(t, tt.expectedUser, userData)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			userStore.AssertExpectations(t)
		})
	}
}

func TestValidateTenantOwnership(t *testing.T) {
	tests := []struct {
		name           string
		tenantID       string
		userID         string
		setupMock      func(*mocks.MockTenant)
		expectedTenant *store.TenantData
		expectedError  error
	}{
		{
			name:     "valid_owner",
			tenantID: "tenant-123",
			userID:   "owner-123",
			setupMock: func(ts *mocks.MockTenant) {
				ts.On("FindOne", mock.Anything, &store.TenantFilter{ID: "tenant-123"}).
					Return(&store.TenantData{
						ID:      "tenant-123",
						OwnerID: sql.NullString{String: "owner-123", Valid: true},
					}, nil)
			},
			expectedTenant: &store.TenantData{
				ID:      "tenant-123",
				OwnerID: sql.NullString{String: "owner-123", Valid: true},
			},
			expectedError: nil,
		},
		{
			name:     "tenant_not_found",
			tenantID: "non-existent",
			userID:   "owner-123",
			setupMock: func(ts *mocks.MockTenant) {
				ts.On("FindOne", mock.Anything, &store.TenantFilter{ID: "non-existent"}).
					Return(nil, sql.ErrNoRows)
			},
			expectedTenant: nil,
			expectedError:  errorHandler.NewNotFound(errorHandler.WithInfo("tenant not found")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenantStore := new(mocks.MockTenant)
			tt.setupMock(tenantStore)

			tenant, err := commonvalidator.ValidateTenantOwnership(context.Background(), tenantStore, tt.tenantID, tt.userID)

			assert.Equal(t, tt.expectedTenant, tenant)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			tenantStore.AssertExpectations(t)
		})
	}
}
