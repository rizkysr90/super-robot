package permission

import (
	"context"
	"fmt"
)

// RuleBasedValidator implements a flexible validation system
type RuleBasedValidator struct {
	rules  []Rule
	stores *ValidationStores
}

// NewRuleBasedValidator creates a new validator with the specified rules
func NewRuleBasedValidator(stores *ValidationStores, rules ...Rule) *RuleBasedValidator {
	return &RuleBasedValidator{
		rules:  rules,
		stores: stores,
	}
}

// First, let's modify our Validate method signature to return both error and the result
func (v *RuleBasedValidator) Validate(ctx context.Context,
	tenantID,
	actionBy,
	branchIDs,
	permissionCode string) (*PermissionResult, error) {
	// Initialize the validation context with the result
	result := &PermissionResult{
		IsAllowed: false,
		IsOwner:   false,
		Data:      &PermissionData{},
	}

	validationCtx := &ValidationContext{
		TenantID:       tenantID,
		ActionBy:       actionBy,
		PermissionCode: permissionCode,
		Result:         result,
		Stores:         v.stores,
		BranchIDs:      branchIDs,
		SkipRemaining:  false,
	}

	// Run each rule in sequence
	for _, rule := range v.rules {
		if err := rule.Validate(ctx, validationCtx); err != nil {
			// Even if validation fails, we return the collected data along with the error
			return result, fmt.Errorf("%s: %w", rule.GetName(), err)
		}
		if validationCtx.SkipRemaining {
			break
		}
	}

	// Return the result if all validations pass
	return result, nil
}
