package dt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

func DataTable[T any](req *Request, db *gorm.DB, model T, preloadRelations ...string) (PaginatedResponse[T], error) {
	var records []T
	var count int64

	// default page to 1
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = 10
	}

	if req.Sort == "" {
		req.Sort = "id"
	}

	if req.Order == "" {
		req.Order = "asc"
	}

	// Membuat query dasar menggunakan koneksi dari struktur
	query := db.Model(&model)

	// Preload relations
	for _, relation := range preloadRelations {
		query = query.Preload(relation)
	}

	// Filter by date range
	if req.DatetimeColumn != "" {
		if req.DatetimeFrom == "" || req.DatetimeTo == "" {
			return PaginatedResponse[T]{}, errors.New("datetime from and to must be provided when datetime column is specified")
		}

		if !isValidColumn(req.DatetimeColumn, model) {
			return PaginatedResponse[T]{}, errors.New("datetime column is not valid")
		}

		var start, end time.Time
		var err error

		start, err = time.Parse(time.RFC3339, req.DatetimeFrom)
		if err != nil {
			return PaginatedResponse[T]{}, fmt.Errorf("invalid datetime from format: %v. use RFC3339 format (ie. 2006-01-02T15:04:05Z07:00)", err)
		}

		end, err = time.Parse(time.RFC3339, req.DatetimeTo)
		if err != nil {
			return PaginatedResponse[T]{}, fmt.Errorf("invalid end date: %v. use RFC3339 format (ie. 2006-01-02T15:04:05Z07:00)", err)
		}

		query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", req.DatetimeColumn), start, end)
	}

	// Search functionality
	if req.SearchValue != "" && req.SearchColumns != "" {
		columns := strings.Split(req.SearchColumns, ",")
		if columns[0] != "" {
			// check column validity
			for i, col := range columns {
				if !isValidColumn(strings.TrimSpace(col), model) {
					return PaginatedResponse[T]{}, errors.New("search column is not valid")
				}

				// if column type not string then convert to string
				if !checkColumnType(col, model, "string") {
					columns[i] = fmt.Sprintf("CAST(%s AS TEXT)", col)
				} else {
					columns[i] = fmt.Sprintf("LOWER(%s)", col)
				}
			}
			searchQuery := strings.Join(columns, fmt.Sprintf(" LIKE '%%%s%%' OR ", strings.ToLower(req.SearchValue)))
			searchQuery = fmt.Sprintf("%s LIKE '%%%s%%'", searchQuery, strings.ToLower(req.SearchValue))
			query = query.Where(searchQuery)
		}
	}

	// Filter functionality
	if req.FilterColumns != "" && req.FilterValue != "" {
		filterColumns := strings.Split(req.FilterColumns, ",")
		filterValues := strings.Split(req.FilterValue, "|")

		if len(filterColumns) != len(filterValues) {
			return PaginatedResponse[T]{}, errors.New("filter columns must equal with filter values")
		}

		for i, filterColumn := range filterColumns {
			// Pisahkan nilai jika ada beberapa value dipisahkan koma
			values := strings.Split(filterValues[i], ",")

			// Trim whitespace dari setiap value
			for j := range values {
				values[j] = strings.TrimSpace(values[j])
			}

			// Cek apakah kolom merupakan tipe teks, jika ya gunakan LOWER()
			if checkColumnType(filterColumn, model, "string") {
				// Gunakan LOWER() pada kolom teks
				query = query.Where(fmt.Sprintf("LOWER(%s) IN (?)", filterColumn), lowerValues(values))
			} else {
				// Jika bukan tipe teks, gunakan query biasa tanpa LOWER()
				query = query.Where(fmt.Sprintf("%s IN (?)", filterColumn), values)
			}
		}
	}

	// Hitung total row
	if err := query.Count(&count).Error; err != nil {
		return PaginatedResponse[T]{}, err
	}

	// Sorting dan Pagination
	if !isValidColumn(req.Sort, model) {
		return PaginatedResponse[T]{}, errors.New("sort column is not valid")

	}
	sortOrder := fmt.Sprintf("%s %s", req.Sort, req.Order)
	offset := (req.Page - 1) * req.PageSize

	result := query.Limit(req.PageSize).Offset(offset).Order(sortOrder).Find(&records)
	if result.Error != nil {
		return PaginatedResponse[T]{}, result.Error
	}

	totalPages := calculateTotalPages(int(count), req.PageSize)

	// Return response struct
	return PaginatedResponse[T]{
		Records:    records,
		TotalCount: int(count),
		TotalPages: totalPages,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Pages:      generatePageArray(totalPages),
	}, nil
}
