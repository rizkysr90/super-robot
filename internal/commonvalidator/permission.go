package commonvalidator

import (
	"context"
	"database/sql"
	"errors"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/pkg/errorHandler"
)

type Permission struct {
	RequiredStore *RequiredStore
	IsAllowed     bool
	IsOwner       bool
	Data          *Data
}
type RequiredStore struct {
	Tenant           store.Tenant
	WorkLocation     store.WorkLocation
	AssignmentRole   store.AssignmentRole
	TenantPermission store.TenantPermission
}
type Data struct {
	Tenant               *store.TenantData
	WorkLocation         []store.WorkLocationData
	AssignmentRole       []store.AssignmentRoleData
	TenantRolePermission []store.TenantPermissionData
}

func NewPermission(
	tenant store.Tenant,
	workLocation store.WorkLocation,
	assignmentRole store.AssignmentRole,
	tenantPermission store.TenantPermission,
) *Permission {
	return &Permission{
		RequiredStore: &RequiredStore{
			Tenant:           tenant,
			WorkLocation:     workLocation,
			AssignmentRole:   assignmentRole,
			TenantPermission: tenantPermission,
		},
		Data: &Data{},
	}
}
func (p *Permission) Validate(ctx context.Context, tenantID, actionBy, permissionCode string) error {
	tenantData, err := p.RequiredStore.Tenant.FindOne(ctx, &store.TenantFilter{ID: tenantID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errorHandler.NewNotFound(errorHandler.WithInfo("IsOwnerValidate : tenant not found"))
		}
		return err
	}
	if tenantData.OwnerID.String == actionBy {
		p.IsOwner = true
		p.IsAllowed = true
		p.Data.Tenant = tenantData
		return nil
	}
	p.Data.Tenant = tenantData
	if !p.IsOwner {
		// do check on roles assignment
		workLocationsData, err := p.RequiredStore.WorkLocation.FindByUserID(ctx, actionBy)
		if err != nil {
			errMsg := "worklocation store : " + err.Error()
			return errors.New(errMsg)
		}
		p.Data.WorkLocation = workLocationsData
		buildAssigmentID := []string{}
		for _, workLocationData := range workLocationsData {
			buildAssigmentID = append(buildAssigmentID, workLocationData.ID)
		}
		assignmentRolesData, err := p.RequiredStore.AssignmentRole.FindMany(ctx, buildAssigmentID)
		if err != nil {
			errMsg := "assignment roles store : " + err.Error()
			return errors.New(errMsg)
		}
		p.Data.AssignmentRole = assignmentRolesData
		buildTenantRoleID := []string{}
		for _, assassignmentRolesData := range assignmentRolesData {
			buildTenantRoleID = append(buildTenantRoleID, assassignmentRolesData.TenantRoleID)
		}
		tenantRolePermissionsData, err := p.RequiredStore.TenantPermission.FindByTenantRoleID(ctx, buildTenantRoleID)
		if err != nil {
			errMsg := "tenant role permission store : " + err.Error()
			return errors.New(errMsg)
		}
		p.Data.TenantRolePermission = tenantRolePermissionsData
		for _, tenantRolePermissionData := range tenantRolePermissionsData {
			if tenantRolePermissionData.PermissionCode == permissionCode {
				p.IsAllowed = true
				p.IsOwner = false
				return nil
			}
		}
	}
	p.IsAllowed = false
	p.IsOwner = false
	return nil
}
