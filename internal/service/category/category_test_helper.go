package category

import (
	store "auth-service-rizkysr90-pos/internal/store/mocks"

	"github.com/DATA-DOG/go-sqlmock"

	"database/sql"
	"testing"
)

// Helper function to initialize the service and mock store
func initTestService(t *testing.T) (
	*Service, *store.MockCategoryStore, sqlmock.Sqlmock, *sql.DB) {
	// Initialize the mock category store
	mockCategoryStore := new(store.MockCategoryStore)
	// Initialize the mock sql.DB and sqlmock
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to initialize sqlmock: %v", err)
	}
	// Initialize the service with the mock store
	svc := &Service{
		db: db,
		categoryStore: mockCategoryStore,
	}

	return svc, mockCategoryStore, sqlMock, db
}