package permission

import (
	"context"
	"database/sql"
	"errors"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/pkg/errorHandler"
)

// OwnerRule checks if the user is the tenant owner
type OwnerRule struct{}

func (r *OwnerRule) GetName() string {
	return "OwnerValidation"
}
func (r *OwnerRule) Validate(ctx context.Context, vc *ValidationContext) error {
	tenantData, err := vc.Stores.TenantFinder.FindOne(ctx, &store.TenantFilter{ID: vc.TenantID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errorHandler.NewNotFound(errorHandler.WithInfo("tenant not found"))
		}
		return err
	}

	vc.Result.Data.Tenant = tenantData
	if tenantData.OwnerID.String == vc.ActionBy {
		vc.Result.IsOwner = true
		vc.Result.IsAllowed = true
		vc.SkipRemaining = true // Skip remaining rules
	}
	return nil
}
