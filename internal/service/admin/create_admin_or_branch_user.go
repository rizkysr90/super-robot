package admin

import (
	"context"
	"database/sql"
	"errors"
	"rizkysr90-pos/internal/commonvalidator"
	"rizkysr90-pos/internal/constant"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/utility"
	"rizkysr90-pos/pkg/errorHandler"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type UserType string
type AuthType string

const (
	Password AuthType = "password"
	Google   AuthType = "google"

	Admin  UserType = "ADMIN"
	Branch UserType = "BRANCH"
)

type reqCreateUser struct {
	*RequestCreateAdmin
}

func (req *reqCreateUser) sanitize() {
	req.TenantID = strings.TrimSpace(req.TenantID)
	req.ActionBy = strings.TrimSpace(req.ActionBy)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.BranchID = strings.TrimSpace(req.BranchID)
	req.FullName = strings.TrimSpace(req.FullName)
}

func (req *reqCreateUser) validate() error {
	validationErrors := []errorHandler.HttpError{}

	// Required field validations
	if !commonvalidator.IsRequired(req.TenantID) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("tenant_id"))
	}
	if !commonvalidator.IsRequired(req.ActionBy) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("action_by"))
	}
	if !commonvalidator.IsRequired(req.Email) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("email"))
	}
	if !commonvalidator.IsRequired(req.Password) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("password"))
	}
	if !commonvalidator.IsRequired(string(req.Type)) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("user_type"))
	}

	// Max length validations
	if !commonvalidator.MaxLen(req.TenantID, 500) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("tenant_id"))
	}
	if !commonvalidator.MaxLen(req.ActionBy, 100) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("action_by"))
	}
	if !commonvalidator.MaxLen(req.Email, 255) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("email"))
	}
	if !commonvalidator.MaxLen(req.Password, 100) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("password"))
	}
	// Branch ID max length validation (only if provided)
	if req.BranchID != "" && !commonvalidator.MaxLen(req.BranchID, 500) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("branch_id"))
	}

	// Email format validation
	if !commonvalidator.IsEmail(req.Email) {
		validationErrors = append(validationErrors, *utility.ConstructErrorInvalid("email"))
	}

	// Password strength validation
	if !commonvalidator.IsStrongPassword(req.Password) {
		validationErrors = append(validationErrors, *utility.ConstructErrorInvalid("password"))
	}

	// User type validation
	if !commonvalidator.IsUserType(string(req.Type)) {
		validationErrors = append(validationErrors, *utility.ConstructErrorInvalid("user_type"))
	}
	// Branch ID is required only for branch users
	if req.Type == Branch && !commonvalidator.IsRequired(req.BranchID) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("branch_id"))
	}
	if len(validationErrors) > 0 {
		return errorHandler.NewMultipleFieldsValidation(validationErrors)
	}
	return nil
}

// RequestCreateAdmin represents the request payload for creating an admin user
// @Description Request body for creating a new admin user
type RequestCreateAdmin struct {
	// Tenant ID for the admin user
	// Required: true
	// Max Length: 500
	// Min Length: 1
	TenantID string `json:"tenant_id" binding:"required,max=500,min=1" example:"550e8400-e29b-41d4-a716-446655440000"`

	// ID of the user performing the action
	// Required: true
	// Max Length: 100
	// Min Length: 1
	ActionBy string `json:"action_by" binding:"required,max=100,min=1" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Email address of the admin user
	// Required: true
	// Max Length: 200
	// Min Length: 3
	// Pattern: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	Email string `json:"email" binding:"required,email,max=200,min=3" example:"admin@example.com"`

	// Password for the admin user
	// Required: true
	// Max Length: 100
	// Min Length: 8
	Password string `json:"password" binding:"required,max=100,min=8" example:"StrongP@ssw0rd"`

	// Type of user (ADMIN or BRANCH)
	// Required: true
	// Enum: ADMIN,BRANCH
	Type UserType `json:"user_type" binding:"required,oneof=ADMIN BRANCH" example:"ADMIN"`

	// Branch ID - UUID of the branch where the user will be assigned
	// Required: true when user_type is BRANCH, optional otherwise
	// Max Length: 500
	// Format: uuid
	// Example: 550e8400-e29b-41d4-a716-446655440000
	// Note: Must be a valid UUID of an existing branch in the system. Required when creating branch users.
	BranchID string `json:"branch_id" binding:"omitempty,uuid,max=500" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Full name
	// Required
	// Max Length : 255
	// Example : Rizki Susilo R
	FullName string `json:"full_name" binding:"max=255" example:"Rizki Susilo Ramadhan"`
}

func (s *Service) CreateUser(ctx context.Context, request *RequestCreateAdmin) error {
	input := &reqCreateUser{request}
	input.sanitize()
	if err := input.validate(); err != nil {
		return err
	}
	// Check duplicate email
	userData, err := s.userStore.FindOne(ctx, &store.UserQueryFilter{
		Email: input.Email,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if userData != nil {
		return errorHandler.NewBadRequest(errorHandler.WithInfo("user duplicate"))
	}
	// Check is branch found
	if input.BranchID != "" {
		_, err := s.branchStore.FindOne(ctx, &store.BranchesFilter{ID: input.BranchID})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errorHandler.NewNotFound(errorHandler.WithInfo("branch not found"))
			}
			return err
		}
	}
	// Check permissions
	permission := commonvalidator.NewPermission(
		s.tenantStore,
		s.worklocationStore,
		s.assignmentRoleStore,
		s.tenantPermissionStore,
		s.userStore,
	)
	permissionCode := constant.RbacNewAdmin
	if input.Type == Branch {
		permissionCode = constant.RbacNewBranchUser
	}
	if err := permission.Validate(
		ctx,
		input.TenantID,
		input.ActionBy,
		permissionCode,
	); err != nil {
		return err
	}
	if !permission.IsAllowed {
		return errorHandler.NewUnauthorized(errorHandler.WithInfo("permission is not allowed"))
	}
	// Set data
	newUserID := uuid.NewString()
	// user data
	hashPassword, err := utility.HashPassword(input.Password)
	if err != nil {
		return err
	}
	insertedUserData := &store.UserData{
		ID:           newUserID,
		Email:        input.Email,
		GoogleID:     sql.NullString{String: "", Valid: false},
		FullName:     input.FullName,
		PasswordHash: sql.NullString{String: hashPassword, Valid: true},
		AuthType:     string(Password),
		UserType:     string(input.Type),
		TenantID:     input.TenantID,
		CreatedAt:    time.Now().UTC(),
		LastLoginAt:  time.Now().UTC(),
		CreatedBy:    sql.NullString{String: input.ActionBy, Valid: true},
	}
	insertedBranchID := sql.NullString{
		String: input.BranchID,
		Valid:  true,
	}
	if insertedBranchID.String == "" {
		insertedBranchID.Valid = false
	}
	insertedAssignmentData := &store.WorkLocationData{
		ID:        uuid.NewString(),
		TenantID:  input.TenantID,
		UserID:    newUserID,
		BranchID:  insertedBranchID,
		CreatedAt: time.Now().UTC(),
		CreatedBy: input.ActionBy,
	}
	err = sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		if err := s.userStore.Insert(txContext, insertedUserData); err != nil {
			return err
		}
		if err := s.worklocationStore.Insert(txContext, insertedAssignmentData); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
