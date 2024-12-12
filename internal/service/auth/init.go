package auth

import (
	"database/sql"
	"rizkysr90-pos/internal/auth"
	"rizkysr90-pos/internal/store"
)

type Auth struct {
	db          *sql.DB
	authClient  *auth.Client
	stateStore  store.State
	userStore   store.User
	tenantStore store.Tenant
	session     store.SessionRedis
}

func NewAuthService(
	sqlDB *sql.DB, authClient *auth.Client, stateStore store.State,
	userStore store.User, tenantStore store.Tenant, session store.SessionRedis,
) *Auth {
	return &Auth{
		db:          sqlDB,
		authClient:  authClient,
		stateStore:  stateStore,
		userStore:   userStore,
		tenantStore: tenantStore,
		session:     session,
	}
}
