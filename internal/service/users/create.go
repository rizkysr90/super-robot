package users

import (
	"auth-service-rizkysr90-pos/internal/helper/errorHandler"
	"auth-service-rizkysr90-pos/internal/helper/validator"
	"auth-service-rizkysr90-pos/internal/store"
	"time"

	payload "auth-service-rizkysr90-pos/internal/payload/http/users"

	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
	"golang.org/x/crypto/bcrypt"
)

type reqCreateUsers struct {
	*payload.ReqCreateUsers
}

func (req *reqCreateUsers) sanitize() {
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.ConfirmPassword = strings.TrimSpace(req.ConfirmPassword)
	req.Phone = strings.TrimSpace(req.Phone)
}
func (req *reqCreateUsers) validateFields() error {
	validationErrors := []errorHandler.HttpError{}
	if !validator.ValidateName(req.FirstName) {
		validationErrors = append(validationErrors,
			*errorHandler.NewBadRequest(errorHandler.WithInfo(
				"invalid format first_name")))
	}
	if !validator.ValidateName(req.LastName) {
		validationErrors = append(validationErrors,
			*errorHandler.NewBadRequest(errorHandler.WithInfo(
				"invalid format lastname")))
	}
	if !validator.ValidateEmail(req.Email) {
		validationErrors = append(validationErrors,
			*errorHandler.NewBadRequest(errorHandler.WithInfo(
				"invalid format email")))
	}
	if !validator.ValidatePassword(req.Password) {
		validationErrors = append(validationErrors, *errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid value password")))
	}
	if req.ConfirmPassword != req.Password {
		validationErrors = append(validationErrors, *errorHandler.NewBadRequest(
			errorHandler.WithInfo("confirm_password not equal to password")))
	}
	if !validator.ValidateOnlyNumber(req.Phone) {
		validationErrors = append(validationErrors, *errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid value phone number")))
	}
	if !validator.ValidateRoles(req.Role) {
		validationErrors = append(validationErrors, *errorHandler.NewBadRequest(
			errorHandler.WithInfo("invalid value roles")))
	}
	if len(validationErrors) > 0 {
		return errorHandler.NewMultipleFieldsValidation(validationErrors)
	}

	return nil
}
func (req *reqCreateUsers) setUserData() (*store.User, error) {
	// gen bcrypt
	bytesPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}
	return &store.User{
		ID:          uuid.NewString(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Password:    string(bytesPassword),
		Email:       req.Email,
		Phone:       req.Phone,
		IsActivated: false,
		Role:        req.Role,
		CreatedAt:   time.Now().UTC(),
	}, nil
}
func (s *Service) CreateUser(ctx context.Context, req *payload.ReqCreateUsers) error {
	input := reqCreateUsers{req}
	input.sanitize()
	err := input.validateFields()
	if err != nil {
		return err
	}
	var getUserByEmail *store.User
	err = sqldb.WithinTx(ctx, s.db, func(qe sqldb.QueryExecutor) error {
		tx := sqldb.WithTxContext(ctx, qe)
		getUserByEmail, err = s.userStore.FindActiveUserByEmail(
			tx, &store.User{Email: req.Email})
		return err
	}, sqldb.TxIsolationLevelReadCommitted())
	if err != nil {
		return err
	}
	if getUserByEmail.ID != "" {
		// users with this email already exist
		return errorHandler.NewBadRequest(
			errorHandler.WithInfo("user with this email has been created"))
	}
	var newUser *store.User
	newUser, err = input.setUserData()
	if err != nil {
		return err
	}
	err = sqldb.WithinTx(ctx, s.db, func(qe sqldb.QueryExecutor) error {
		tx := sqldb.WithTxContext(ctx, qe)
		getUserByEmail, err = s.userStore.FindActiveUserByEmail(
			tx, &store.User{Email: req.Email})
		if err != nil {
			return err
		}
		if getUserByEmail.ID != "" {
			// users with this email already exist
			return errorHandler.NewBadRequest(
				errorHandler.WithInfo("user with this email has been created"))
		}
		return s.userStore.Create(tx, newUser)

	}, sqldb.TxIsolationLevelReadCommitted())
	if err != nil {
		return err
	}
	return nil
}
