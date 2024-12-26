package dt

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mock model struct for testing
type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Test the DataTable function
func TestDataTable(t *testing.T) {
	// Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	// Migrate the schema
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Seed some data
	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
		{ID: 3, Name: "Alice Smith", Email: "alice@example.com"},
	}
	db.Create(&users)

	// Test case 1: Simple pagination and sorting
	req := &Request{
		Page:     1,
		PageSize: 2,
		Sort:     "id",
		Order:    "asc",
	}

	// Run DataTable function
	resp, err := DataTable(req, db, User{})
	if err != nil {
		t.Errorf("DataTable() failed: %v", err)
	}

	// Check the result
	if resp.TotalCount != 3 {
		t.Errorf("Expected TotalCount = 3, got %d", resp.TotalCount)
	}
	if len(resp.Records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(resp.Records))
	}
	if resp.Records[0].Name != "John Doe" {
		t.Errorf("Expected first record to be John Doe, got %s", resp.Records[0].Name)
	}
	if resp.TotalPages != 2 {
		t.Errorf("Expected TotalPages = 2, got %d", resp.TotalPages)
	}

	// Add more test cases for search, filter, and invalid cases as needed
}
