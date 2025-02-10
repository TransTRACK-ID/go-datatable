package dt

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mock model struct for testing
type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// DataTableTestSuite defines the test suite for DataTable
type DataTableTestSuite struct {
	suite.Suite
	DB *gorm.DB
}

// SetupSuite sets up the test database
func (suite *DataTableTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		suite.T().Fatalf("failed to migrate database: %v", err)
	}

	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
		{ID: 3, Name: "Alice Smith", Email: "alice@example.com"},
	}
	db.Create(&users)

	suite.DB = db
}

// TestSimplePaginationAndSorting tests simple pagination and sorting
func (suite *DataTableTestSuite) TestSimplePaginationAndSorting() {
	req := &Request{
		Page:     1,
		PageSize: 2,
		Sort:     "id",
		Order:    "asc",
	}

	resp, err := DataTable(req, suite.DB, User{})
	suite.NoError(err)
	suite.Equal(3, resp.TotalCount)
	suite.Len(resp.Records, 2)
	suite.Equal("John Doe", resp.Records[0].Name)
	suite.Equal(2, resp.TotalPages)
}

// TestSearchFunctionality tests search functionality case-insensitive
func (suite *DataTableTestSuite) TestSearchFunctionality() {
	req := &Request{
		Page:          1,
		PageSize:      2,
		Sort:          "id",
		Order:         "asc",
		SearchColumns: "name, email",
		SearchValue:   "doe",
	}

	resp, err := DataTable(req, suite.DB, User{})
	suite.NoError(err)
	suite.Equal(2, resp.TotalCount)
	suite.Len(resp.Records, 2)
}

// TestInvalidPageNumber tests invalid page number
func (suite *DataTableTestSuite) TestInvalidPageNumber() {
	req := &Request{
		Page:     2,
		PageSize: 2,
		Sort:     "id",
		Order:    "asc",
	}

	resp, err := DataTable(req, suite.DB, User{})
	suite.NoError(err)
	suite.Equal(3, resp.TotalCount)
	suite.Len(resp.Records, 1)
	suite.Equal("Alice Smith", resp.Records[0].Name)
	suite.Equal(2, resp.TotalPages)
}

// TestFilteringFunctionality tests filtering functionality
func (suite *DataTableTestSuite) TestFilteringFunctionality() {
	req := &Request{
		Page:          1,
		PageSize:      2,
		Sort:          "id",
		Order:         "asc",
		FilterColumns: "email",
		FilterValue:   "alice@example.com",
	}

	resp, err := DataTable(req, suite.DB, User{})
	suite.NoError(err)
	suite.Equal(1, resp.TotalCount)
	suite.Len(resp.Records, 1)
	suite.Equal("Alice Smith", resp.Records[0].Name)
}

// TestInvalidSortField tests invalid sort field
func (suite *DataTableTestSuite) TestInvalidSortField() {
	req := &Request{
		Page:     1,
		PageSize: 2,
		Sort:     "invalid_field",
		Order:    "asc",
	}

	_, err := DataTable(req, suite.DB, User{})
	suite.Error(err)
}

// TestDataTable runs the test suite
func TestDataTable(t *testing.T) {
	suite.Run(t, new(DataTableTestSuite))
}
