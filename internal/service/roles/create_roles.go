package roles

import (
	"context"
	"rizkysr90-pos/internal/commonvalidator"
	"rizkysr90-pos/internal/constant"
	"rizkysr90-pos/internal/utility"
	"rizkysr90-pos/pkg/errorHandler"
	"strings"
)

type RequestCreateRoles struct {
	// Tenant ID for the role
	// Required: true
	// Max Length: 500
	TenantID string `json:"tenant_id"`

	// ID of the user performing the action
	// Required: true
	// Max Length: 100
	ActionBy string `json:"action_by"`

	// Name of the role
	// Required: true
	// Max Length: 200
	Name string `json:"name"`

	// Indicates if this is a head office role
	// Required: true
	// Default: false
	IsHeadOfficeRole bool `json:"is_head_office_role"`

	// Description of the role and its responsibilities
	// Required: true
	// Max Length: 500
	Description string `json:"description"`
}
type reqCreateRoles struct {
	*RequestCreateRoles
}

func (req *reqCreateRoles) sanitize() {
	req.TenantID = strings.TrimSpace(req.TenantID)
	req.ActionBy = strings.TrimSpace(req.ActionBy)
	req.Name = strings.TrimSpace(strings.ToUpper(req.Name))
	req.Description = strings.TrimSpace(strings.ToUpper(req.Description))
}

func (req *reqCreateRoles) validate() error {
	validationErrors := []errorHandler.HttpError{}

	// Required field validations
	if !commonvalidator.IsRequired(req.TenantID) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("tenant_id"))
	}
	if !commonvalidator.IsRequired(req.ActionBy) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("action_by"))
	}
	if !commonvalidator.IsRequired(req.Name) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("name"))
	}
	if !commonvalidator.IsRequired(req.Description) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("description"))
	}

	// Max length validations
	if !commonvalidator.MaxLen(req.TenantID, 500) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("tenant_id"))
	}
	if !commonvalidator.MaxLen(req.ActionBy, 100) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("action_by"))
	}
	if !commonvalidator.MaxLen(req.Name, 200) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("name"))
	}
	if !commonvalidator.MaxLen(req.Description, 500) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("description"))
	}

	if len(validationErrors) > 0 {
		return errorHandler.NewMultipleFieldsValidation(validationErrors)
	}
	return nil
}
func (s *Service) Create(ctx context.Context, request *RequestCreateRoles) error {
	input := &reqCreateRoles{request}
	input.sanitize()
	if err := input.validate(); err != nil {
		return err
	}
	permission := commonvalidator.NewPermission(
		s.tenantStore,
		s.worklocationStore,
		s.assignmentRoleStore,
		s.tenantPermissionStore)

	if err := permission.Validate(
		ctx,
		input.TenantID,
		input.ActionBy,
		// "TESTING",
		constant.RbacNewRole,
	); err != nil {
		return err
	}
	if !permission.IsAllowed {
		return errorHandler.NewUnauthorized()
	}
	return nil
}
