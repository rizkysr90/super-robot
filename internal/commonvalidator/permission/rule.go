package permission

import (
	"context"
	"rizkysr90-pos/internal/store"
)

// Rule defines the contract for all validation rules
// Each rule can access the validation context and make its own decision
type Rule interface {
	// Validate performs the rule's validation logic
	Validate(ctx context.Context, vc *ValidationContext) error

	// GetName returns the rule's name for identification and error reporting
	GetName() string
}

// ValidationContext holds all the data needed during validation
// This acts as a shared state that rules can read from and write to
type ValidationContext struct {
	// Input parameters
	TenantID       string
	ActionBy       string
	PermissionCode string

	// Validation results and collected data
	Result *PermissionResult

	// Access to external services
	Stores *ValidationStores

	SkipRemaining bool // New flag

}

// ValidationStores provides access to all required data stores
type ValidationStores struct {
	UserFinder           store.UserFinder
	TenantFinder         store.TenantFinder
	WorkLocationFinder   store.WorkLocationFinder
	AssignmentRoleFinder store.AssignmentRoleFinder
	PermissionFinder     store.TenantPermissionFinder
}

type PermissionResult struct {
	IsAllowed bool
	IsOwner   bool
	Data      *PermissionData
}

type PermissionData struct {
	Tenant               *store.TenantData
	WorkLocation         []store.WorkLocationData
	AssignmentRole       []store.AssignmentRoleData
	TenantRolePermission []store.TenantPermissionData
}
