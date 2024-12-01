package users

import (
	"database/sql"
	"rizkysr90-pos/internal/store"
	jwttoken "rizkysr90-pos/pkg/jwt"
)

type UserService struct {
	db        *sql.DB
	userStore store.UserStore
	jwt       *jwttoken.JWT
}

func NewUsersService(sqlDB *sql.DB, userStore store.UserStore, jwt *jwttoken.JWT) *UserService {
	return &UserService{
		db:        sqlDB,
		userStore: userStore,
		jwt:       jwt,
	}
}
