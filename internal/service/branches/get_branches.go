package branches

import "rizkysr90-pos/internal/store"

type RequestGetBranches struct {
	// Pagination page size
	// Required: false
	// Default: 20
	// minimum: 1
	// maximum: 100
	PageSize int `json:"page_size"`

	// Pagination page number
	// Required: false
	// Default: 1
	// minimum: 1
	PageNumber int `json:"page_number"`

	// Branch name filter
	// Required: false
	// Max Length: 100
	// Pattern: ^[a-zA-Z0-9-_]+$
	BranchName string `json:"branch_name"`

	// Tenant identifier
	// Required: false
	// Required when branch_name is not empty
	// Max Length: 500
	TenantID string `json:"tenant_id" binding:"required_with=BranchName"`

	// User performing the action
	// Required: true
	// Max Length: 500
	ActionBy string `json:"action_by" binding:"required"`
}

// BranchesData represents a branch entity
type BranchesData struct {
	// Branch UUID
	// Required: true
	// Example: 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id"`

	// Tenant UUID
	// Required: true
	// Example: 550e8400-e29b-41d4-a716-446655440001
	TenantID string `json:"tenant_id"`

	// Branch name
	// Required: true
	// Max Length: 255
	// Example: Main Branch Jakarta
	Name string `json:"name"`

	// Branch address
	// Required: false
	// Example: Jl. Sudirman No. 123, Jakarta
	Address string `json:"address,omitempty"`

	// Creation timestamp
	// Required: true
	// Example: 2025-01-11T14:30:00Z
	CreatedAt string `json:"created_at"`

	// Owner who created the branch
	// Required: false
	// Max Length: 255
	// Example: John Doe
	CreatedByOwner string `json:"created_by_owner,omitempty"`

	// User who created the branch
	// Required: false
	// Max Length: 255
	// Example: jane.doe@company.com
	CreatedBy string `json:"created_by,omitempty"`
}
type ResponseGetBranches struct {
	Pagination store.Pagination
	Data       []BranchesData
}
