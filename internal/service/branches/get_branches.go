package branches

import (
	"context"
	"rizkysr90-pos/internal/commonvalidator"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/utility"
	"rizkysr90-pos/pkg/errorHandler"
	"strings"
)

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

type reqGetBranches struct {
	*RequestGetBranches
}

func (req *reqGetBranches) sanitize() {
	// Convert branch name to uppercase and trim spaces if it exists
	if req.BranchName != "" {
		req.BranchName = strings.ToUpper(strings.TrimSpace(req.BranchName))
	}

	// Trim spaces from tenant ID if it exists
	if req.TenantID != "" {
		req.TenantID = strings.TrimSpace(req.TenantID)
	}

	// Trim spaces from action by
	if req.ActionBy != "" {
		req.ActionBy = strings.TrimSpace(req.ActionBy)
	}
}
func (req *reqGetBranches) validate() error {
	validationErrors := []errorHandler.HttpError{}
	if !commonvalidator.MaxLen(req.BranchName, 40) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("branch_name"))
	}
	if !commonvalidator.MaxLen(req.TenantID, 500) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("tenant_id"))
	}
	if !commonvalidator.MaxLen(req.ActionBy, 100) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("created_by"))
	}
	if len(validationErrors) > 0 {
		return errorHandler.NewMultipleFieldsValidation(validationErrors)
	}
	return nil
}
func (s *Service) GetBranches(ctx context.Context, request *RequestGetBranches) (*ResponseGetBranches, error) {
	input := &reqGetBranches{request}
	input.sanitize()
	if err := input.validate(); err != nil {
		return nil, err
	}
	_, _, err := commonvalidator.IsAllowedUser(ctx, s.userStore, request.ActionBy)
	if err != nil {
		return nil, err
	}
	_, err = commonvalidator.ValidateTenantOwnership(
		ctx, s.tenantStore, input.TenantID, input.ActionBy)
	if err != nil {
		return nil, err
	}
	branchesData, pagination, err := s.branchStore.FindManyWithPaginated(ctx, &store.BranchesFilter{
		Name:       request.BranchName,
		TenantID:   request.TenantID,
		PageSize:   request.PageSize,
		PageNumber: request.PageNumber,
	})
	if err != nil {
		return nil, err
	}
	result := &ResponseGetBranches{}
	for _, branch := range branchesData {
		data := &BranchesData{
			ID:             branch.BranchID,
			TenantID:       branch.TenantID,
			Name:           branch.Name,
			Address:        branch.Address,
			CreatedAt:      branch.CreatedBy.String,
			CreatedByOwner: branch.CreatedByOwner.String,
			CreatedBy:      branch.CreatedBy.String,
		}
		result.Data = append(result.Data, *data)
	}
	return &ResponseGetBranches{
		Pagination: *pagination,
		Data:       []BranchesData{},
	}, nil
}
