package store

import "context"

type UserFinder interface {
	FindOne(ctx context.Context, filter *UserQueryFilter) (*UserData, error)
}

type TenantFinder interface {
	FindOne(ctx context.Context, filter *TenantFilter) (*TenantData, error)
}

type WorkLocationFinder interface {
	FindByUserID(ctx context.Context, userID string) ([]WorkLocationData, error)
}

type AssignmentRoleFinder interface {
	FindMany(ctx context.Context, ids []string) ([]AssignmentRoleData, error)
}

type TenantPermissionFinder interface {
	FindByTenantRoleID(ctx context.Context, roleIDs []string) ([]TenantPermissionData, error)
}
