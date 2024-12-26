package dt

type Request struct {
	Page          int    `json:"page" form:"page" query:"page"`
	PageSize      int    `json:"page_size" form:"page_size" query:"page_size"`
	SearchColumns string `json:"search_columns" form:"search_columns" query:"search_columns"`
	SearchValue   string `json:"search_value" form:"search_value" query:"search_value"`
	FilterColumns string `json:"filter_columns" form:"filter_columns" query:"filter_columns"`
	FilterValue   string `json:"filter_value" form:"filter_value" query:"filter_value"`
	Sort          string `json:"sort" form:"sort" query:"sort"`
	Order         string `json:"order" form:"order" query:"order"`
}

type PaginatedResponse[T any] struct {
	Records    []T   `json:"records"`
	TotalCount int   `json:"totalCount"`
	TotalPages int   `json:"totalPages"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Pages      []int `json:"pages"`
}
