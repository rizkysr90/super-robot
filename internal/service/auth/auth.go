package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	payload "github.com/rizkysr90/go-boilerplate/internal/payload/http/auth"
	"github.com/rizkysr90/go-boilerplate/internal/store"
	"github.com/rizkysr90/go-boilerplate/pkg/restapierror"
	"github.com/rizkysr90/go-boilerplate/pkg/sqldb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	db        *sql.DB
	userStore store.UserStore
}

func NewAuthService(db *sql.DB, userStore store.UserStore) *Service {
	return &Service{
		db:        db,
		userStore: userStore,
	}
}

func (s *Service) CreateUser(ctx context.Context, req *payload.ReqCreateAccount) error {
	if req.Password != req.ConfirmPassword {
		return restapierror.NewBadRequest(restapierror.WithMessage("failed to confirm password"))
	}
	// check to db
	filterBy := store.UserFilterBy{
		Email: req.Email,
	}
	result, err := s.userStore.FindOne(ctx, &filterBy, "findactiveuser")
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		// internal server error after call db
		return err
	}
	// check duplicate email
	if result != nil {
		return restapierror.NewBadRequest(restapierror.WithMessage("email already registered"))
	}
	// Create user
	insertedData := store.InsertedData{
		ID:        uuid.NewString(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: time.Now().UTC(),
	}
	if err = sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return s.userStore.Create(txContext, &insertedData)
	}); err != nil {
		return err
	}

	return nil
}
