package dt

import (
	"math"
	"reflect"
	"strings"
)

// isValidColumn are checker for column are valid attribute of model
func isValidColumn(col string, model interface{}) bool {
	// Get the type of the struct
	t := reflect.TypeOf(model)

	// Iterate over the struct fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check if the field name matches the provided column name
		if field.Name == col {
			return true
		}

		// Check if the struct tag (e.g., json tag) matches the provided column name
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && strings.Split(jsonTag, ",")[0] == col {
			return true
		}
	}
	return false
}

// checkColumnType are checker for column type in model
func checkColumnType(col string, model interface{}, colType string) bool {
	// Get the type of the struct
	t := reflect.TypeOf(model)

	// Iterate over the struct fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check if the field name matches the provided column name
		if field.Name == col {
			return field.Type.Name() == colType
		}

		// Check if the struct tag (e.g., json tag) matches the provided column name
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && strings.Split(jsonTag, ",")[0] == col {
			return field.Type.Name() == colType
		}
	}
	return false
}

// lowerValues are convert values to lower string
func lowerValues(values []string) []string {
	for i, v := range values {
		values[i] = strings.ToLower(v)
	}
	return values
}

// generatePageArray are array of pages start from 1
func generatePageArray(totalPages int) []int {
	pages := make([]int, totalPages)
	for i := 0; i < totalPages; i++ {
		pages[i] = i + 1
	}
	return pages
}

// calculateTotalPages is a function to calculate  total pages based on total count divide by page size
func calculateTotalPages(totalCount, pageSize int) int {
	if pageSize == 0 {
		// Jika pageSize adalah 0, berarti tidak ada limit, dan semua data ada di satu halaman
		return 1
	}
	return int(math.Ceil(float64(totalCount) / float64(pageSize)))
}
