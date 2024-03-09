package auth

import (
	"context"
	"database/sql"
	"strings"
	"time"

	payload "github.com/rizkysr90/go-boilerplate/internal/payload/http/auth"
	"github.com/rizkysr90/go-boilerplate/internal/store"
	jwttoken "github.com/rizkysr90/go-boilerplate/pkg/jwt"
	"github.com/rizkysr90/go-boilerplate/pkg/restapierror"
	"github.com/rizkysr90/go-boilerplate/pkg/sqldb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db        *sql.DB
	userStore store.UserStore
	jwtToken  jwttoken.JWT
}

func NewAuthService(db *sql.DB, userStore store.UserStore, jwttoken jwttoken.JWT) *Service {
	return &Service{
		db:        db,
		userStore: userStore,
		jwtToken:  jwttoken,
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
	if err != nil {
		if !strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return err
		}
	}
	// check duplicate email
	if result != nil {
		return restapierror.NewBadRequest(restapierror.WithMessage("email already registered"))
	}
	// gen bcrypt
	bytesPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return err
	}
	// Create user
	insertedData := store.InsertedData{
		ID:          uuid.NewString(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    string(bytesPassword),
		CreatedAt:   time.Now().UTC(),
		IsActivated: true,
	}
	if err = sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return s.userStore.Create(txContext, &insertedData)
	}); err != nil {
		return err
	}

	return nil
}
func (s *Service) LoginUser(ctx context.Context,
	req *payload.ReqLoginUser) (*payload.ResLoginUser, error) {
	result, err := s.userStore.FindOne(ctx, &store.UserFilterBy{
		Email: req.Email,
	}, "findactiveuser")
	if err != nil {
		if strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return nil, restapierror.NewNotFound(restapierror.WithMessage("user not found"))
		}
	}
	if err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(req.Password)); err != nil {
		return nil, restapierror.NewBadRequest(restapierror.WithMessage(err.Error()))
	}
	// gen token
	var genToken string
	genToken, err = s.jwtToken.Generate(&jwttoken.JWTClaims{
		UserID: result.ID,
	})
	if err != nil {
		return nil, err
	}

	return &payload.ResLoginUser{
		Token: genToken,
	}, nil
}
