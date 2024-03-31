package store

type Pagination struct {
	PageSize      int `json:"page_size"`
	PageNumber    int `json:"page_number"`
	TotalPages    int `json:"total_pages"`
	TotalElements int `json:"total_elements"`
}
