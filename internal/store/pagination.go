package store

// Pagination represents the pagination information
type Pagination struct {
	// Current page size
	// Required: true
	// Example: 20
	// Minimum: 1
	PageSize int `json:"page_size"`

	// Current page number
	// Required: true
	// Example: 1
	// Minimum: 1
	PageNumber int `json:"page_number"`

	// Total number of pages
	// Required: true
	// Example: 5
	// Minimum: 0
	TotalPages int `json:"total_pages"`

	// Total number of elements
	// Required: true
	// Example: 100
	// Minimum: 0
	TotalElements int `json:"total_elements"`
}
