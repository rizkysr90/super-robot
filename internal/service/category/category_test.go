package category

import (
	store "rizkysr90-pos/internal/store/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

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
		db:            db,
		categoryStore: mockCategoryStore,
	}

	return svc, mockCategoryStore, sqlMock, db
}
func TestNewCategoryService(t *testing.T) {
	// Create a mock SQL database connection (using a real DB is not necessary for this test)
	db := &sql.DB{}

	// Create a mock category store
	mockCategoryStore := new(store.MockCategoryStore)

	// Initialize the service
	svc := NewCategoryService(db, mockCategoryStore)

	// Assert that the service is initialized with the correct values
	assert.Equal(t, db, svc.db)
	assert.Equal(t, mockCategoryStore, svc.categoryStore)
}
