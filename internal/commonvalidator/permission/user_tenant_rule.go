package permission

import (
	"context"
	"database/sql"
	"errors"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/pkg/errorHandler"
)

// UserTenantRule validates that the user belongs to the specified tenant
type UserTenantRule struct{}

func (r *UserTenantRule) GetName() string {
	return "UserTenantValidation"
}
func (r *UserTenantRule) Validate(ctx context.Context, vc *ValidationContext) error {
	userData, err := vc.Stores.UserFinder.FindOne(ctx, &store.UserQueryFilter{ID: vc.ActionBy})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errorHandler.NewNotFound(errorHandler.WithInfo("user not found"))
		}
		return err
	}

	if userData.TenantID != vc.TenantID {
		return errorHandler.NewUnauthorized(errorHandler.WithInfo("invalid access"))
	}
	return nil
}
