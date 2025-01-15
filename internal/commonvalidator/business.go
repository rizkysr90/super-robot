package commonvalidator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/pkg/errorHandler"
)

type IsOwnerData struct {
	Tenant *store.TenantData
	Check  bool
}

func IsOwner(ctx context.Context, tenantStore store.Tenant, tenantID, actionBy string) (*IsOwnerData, error) {
	result := &IsOwnerData{}
	tenantData, err := tenantStore.FindOne(ctx, &store.TenantFilter{ID: tenantID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorHandler.NewNotFound(errorHandler.WithInfo("IsOwnerValidate : tenant not found"))
		}
		return nil, err
	}
	if tenantData.OwnerID.String == actionBy {
		result.Check = true
	}
	result.Tenant = tenantData
	return result, nil
}
func IsAllowedUser(ctx context.Context, userStore store.User, actionByUserID string) (bool, *store.UserData, error) {
	// First check if user exists
	userData, err := userStore.FindOne(ctx, &store.UserQueryFilter{ID: actionByUserID})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, nil, fmt.Errorf("failed to get user data: %w", err)
	}

	isAllowed := false

	// Check if user exists in main user table
	if userData != nil {
		return true, userData, nil
	}

	// Add additional checks here for other tables
	// Example:
	// adminData, err := adminStore.FindOne(ctx, actionByUserID)
	// if err != nil && !errors.Is(err, sql.ErrNoRows) {
	//     return false, nil, fmt.Errorf("failed to check admin permissions: %w", err)
	// }
	// if adminData != nil {
	//     isAllowed = true
	//     return true, userData, nil
	// }

	// If no permissions found in any table
	if !isAllowed {
		return false, nil, errorHandler.NewUnauthorized(errorHandler.WithInfo("not allowed"))
	}

	return false, userData, nil
}

func ValidateTenantOwnership(ctx context.Context,
	tenantStore store.Tenant, tenantID, actionByUserID string) (*store.TenantData, error) {
	existingTenant, err := tenantStore.FindOne(ctx, &store.TenantFilter{ID: tenantID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorHandler.NewNotFound(errorHandler.WithInfo("tenant not found"))
		}
		return nil, err
	}
	if existingTenant.OwnerID.String != actionByUserID {
		return nil, errorHandler.NewBadRequest(errorHandler.WithInfo("invalid user request"))
	}
	return existingTenant, nil
}
