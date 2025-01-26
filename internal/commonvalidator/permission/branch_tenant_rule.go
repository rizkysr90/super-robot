package permission

import (
	"context"
	"errors"
	"strings"
)

// BranchTenantRule checks if the branchID is in tenantID
type BranchTenantRule struct{}

func (r *BranchTenantRule) GetName() string {
	return "BranchTenantValidation"
}

func (r *BranchTenantRule) Validate(ctx context.Context, vc *ValidationContext) error {
	if vc.Stores.BranchFinder == nil {
		return errors.New("branch finder interface implementation is nil")
	}
	if vc.BranchIDs == "" {
		return nil
	}
	branchIDsParams := strings.Split(vc.BranchIDs, ",")
	branchNotFound := make(map[string]int)
	for _, branchID := range branchIDsParams {
		_, exist := branchNotFound[branchID]
		if exist {
			return errors.New("branch id duplicated")
		}
		branchNotFound[branchID] = 1
	}
	branchesData, err := vc.Stores.BranchFinder.FindMany(ctx, branchIDsParams)
	if err != nil {
		return err
	}
	for _, branch := range branchesData {
		_, exist := branchNotFound[branch.BranchID]
		if exist {
			delete(branchNotFound, branch.BranchID)
		}
		if branch.TenantID != vc.TenantID {
			return errors.New("invalid branch id and tenant id")
		}
	}
	if len(branchNotFound) != 0 {
		return errors.New("some branches not found")
	}
	vc.Result.Data.Branches = branchesData
	return nil
}
