# Datatable Package for Golang

This Go package provides a flexible and efficient way to handle paginated data tables, including support for sorting, filtering, and searching using [GORM](https://gorm.io/). It's designed to simplify building API endpoints that return paginated results from a database.

# Requirement
- [GORM](https://gorm.io/) ORM library for Golang

## Features
- **Pagination**: Automatically handles paginated responses with customizable page size.
- **Sorting**: Allows sorting by any column in ascending or descending order.
- **Search**: Provides search functionality across multiple columns.
- **Filtering**: Supports filtering by specific columns with flexible filter values.
- **Preload Relations**: Preload related models in GORM queries.

## Installation

To install the package, run this command via CLI inside your project:

```bash
go get github.com/TransTRACK-ID/go-datatable
```

## Usage
1. Setup your model and request structs
You need to define your model, request, and response structures.

```
// User model
type User struct {
    ID    uint      `json:"id"`
    Name  string    `json:"name"`
    Email string    `json:"email"`
}

// Request struct
type Request struct {
    Page         int
    PageSize     int
    Sort         string
    Order        string
    SearchValue  string
    SearchColumns string
    FilterColumns string
    FilterValue   string
}
```

2. Call the Datatable function
```
func main() {
    db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

    // Example request
    req := &datatable.Request{
        Page:         1,
        PageSize:     10,
        Sort:         "id",
        Order:        "asc",
        SearchValue:  "John",
        SearchColumns: "name,email",
        FilterColumns: "email",
        FilterValue:   "john@example.com",
    }

    // Call the DataTable function
    response, err := dt.DataTable(req, db, User{})
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Print the paginated response
    fmt.Printf("Total Count: %d, Total Pages: %d\n", response.TotalCount, response.TotalPages)
    for _, user := range response.Records {
        fmt.Printf("User: %s, Email: %s\n", user.Name, user.Email)
    }
}
```

example result
```
Total Count: 1, Total Pages: 1
User: John Doe, Email: john@example.com
```