package permission

import (
	"context"
	"fmt"
)

type RolePermissionRule struct{}

func (r *RolePermissionRule) GetName() string {
	return "RolePermissionValidation"
}

func (r *RolePermissionRule) Validate(ctx context.Context, vc *ValidationContext) error {
	if vc.Result.IsOwner {
		return nil
	}

	workLocations, err := vc.Stores.WorkLocationFinder.FindByUserID(ctx, vc.ActionBy)
	if err != nil {
		return fmt.Errorf("worklocation store: %w", err)
	}
	vc.Result.Data.WorkLocation = workLocations

	var assignmentIDs []string
	for _, loc := range workLocations {
		assignmentIDs = append(assignmentIDs, loc.ID)
	}

	assignments, err := vc.Stores.AssignmentRoleFinder.FindMany(ctx, assignmentIDs)
	if err != nil {
		return fmt.Errorf("assignment roles store: %w", err)
	}
	vc.Result.Data.AssignmentRole = assignments

	var roleIDs []string
	for _, assignment := range assignments {
		roleIDs = append(roleIDs, assignment.TenantRoleID)
	}

	permissions, err := vc.Stores.PermissionFinder.FindByTenantRoleID(ctx, roleIDs)
	if err != nil {
		return fmt.Errorf("tenant role permission store: %w", err)
	}
	vc.Result.Data.TenantRolePermission = permissions

	for _, perm := range permissions {
		if perm.PermissionCode == vc.PermissionCode {
			vc.Result.IsAllowed = true
			return nil
		}
	}

	return nil
}
