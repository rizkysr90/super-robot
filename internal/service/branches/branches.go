package branches

import (
	"database/sql"
	"rizkysr90-pos/internal/config"
	"rizkysr90-pos/internal/store"
)

type Service struct {
	db          *sql.DB
	cfg         *config.Config
	tenantStore store.Tenant
	userStore   store.User
	branchStore store.Branches
}

func NewBranchService(sqlDB *sql.DB, cfg *config.Config, tenantStore store.Tenant,
	userStore store.User, branchStore store.Branches) *Service {
	return &Service{
		db:          sqlDB,
		cfg:         cfg,
		tenantStore: tenantStore,
		userStore:   userStore,
		branchStore: branchStore,
	}
}
