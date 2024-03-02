package auth

import (
	payload "api-iad-ams/internal/payload/http/auth"
	"api-iad-ams/internal/store"
	"api-iad-ams/pkg/restapierror"
	"api-iad-ams/pkg/sqldb"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	db        *pgxpool.Pool
	userStore store.UserStore
}

func NewAuthService(db *pgxpool.Pool, userStore store.UserStore) *AuthService {
	return &AuthService{
		db:        db,
		userStore: userStore,
	}
}

func (a *AuthService) CreateUser(ctx context.Context, req *payload.ReqCreateAccount) error {
	if req.Password != req.ConfirmPassword {
		return restapierror.NewBadRequest(ctx, restapierror.WithMessage("failed to confirm password"))
	}
	// check to db
	filterBy := store.UserData{
		Email: req.Email,
	}
	result, err := a.userStore.FindOne(ctx, &filterBy, "findbyemail")
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		// internal server error after call db
		return err
	}
	// check duplicate email
	if result != nil {
		return restapierror.NewBadRequest(ctx, restapierror.WithMessage("email already registered"))
	}
	// Create user
	createdData := store.UserData{
		Id:        uuid.NewString(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: time.Now().UTC(),
	}
	if err = sqldb.WithinTx(ctx, a.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return a.userStore.Create(txContext, &createdData)
	}); err != nil {
		return err
	}

	return nil
}
