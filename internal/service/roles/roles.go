package roles

import (
	"database/sql"
	"rizkysr90-pos/internal/config"
	"rizkysr90-pos/internal/store"
)

type Service struct {
	db                    *sql.DB
	cfg                   *config.Config
	tenantStore           store.Tenant
	userStore             store.User
	branchStore           store.Branches
	assignmentRoleStore   store.AssignmentRole
	tenantPermissionStore store.TenantPermission
	worklocationStore     store.WorkLocation
}

func NewService(
	sqlDB *sql.DB,
	cfg *config.Config,
	tenantStore store.Tenant,
	userStore store.User,
	branchStore store.Branches,
	assignmentRoleStore store.AssignmentRole,
	tenantPermissionStore store.TenantPermission,
	worklocationStore store.WorkLocation,
) *Service {
	return &Service{
		db:                    sqlDB,
		cfg:                   cfg,
		tenantStore:           tenantStore,
		userStore:             userStore,
		branchStore:           branchStore,
		assignmentRoleStore:   assignmentRoleStore,
		tenantPermissionStore: tenantPermissionStore,
		worklocationStore:     worklocationStore,
	}
}
