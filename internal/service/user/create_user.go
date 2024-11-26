package users

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/internal/store"
	"context"
	"database/sql"
	"errors"
	"time"

	"auth-service-rizkysr90-pos/pkg/errorHandler"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserService) CreateUser(ctx context.Context, req *payload.ReqCreateUsers) (*payload.ResCreateUsers, error) {
	if req.ConfirmPassword != req.Password {
		return nil, errorHandler.NewBadRequest(errorHandler.WithMessage("invalid password"))
	}
	user, err := u.userStore.FindByUsername(ctx, req.Username)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return nil, err
	}
	if user != nil {
		return nil, errorHandler.NewBadRequest(errorHandler.WithMessage("duplicate username"))
	}
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}
	insertedUserData := &store.User{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Username:  req.Username,
		Password:  string(password),
	}
	err = sqldb.WithinTx(ctx, u.db, func(qe sqldb.QueryExecutor) error {
		tx := sqldb.WithTxContext(ctx, qe)
		return u.userStore.Create(tx, insertedUserData)
	})
	if err != nil {
		return nil, err
	}
	return &payload.ResCreateUsers{}, nil
}
