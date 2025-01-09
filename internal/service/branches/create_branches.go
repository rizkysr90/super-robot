package branches

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rizkysr90-pos/internal/commonvalidator"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/utility"
	"rizkysr90-pos/pkg/errorHandler"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type RequestCreate struct {
	// Tenant ID for the branch
	// Required: true
	// Max Length: 500
	TenantID string `json:"tenant_id"`

	// Name of the branch
	// Required: true
	// Max Length: 40
	BranchName string `json:"branch_name"`

	// Physical address of the branch
	// Required: true
	// Max Length: 500
	Address string `json:"address"`

	// ID of the user creating the branch
	// Required: true
	// Max Length: 100
	CreatedBy string `json:"created_by"`
}

type reqCreate struct {
	*RequestCreate
}

func (req *reqCreate) sanitize() {
	req.TenantID = strings.TrimSpace(req.TenantID)
	req.BranchName = strings.TrimSpace(strings.ToUpper(req.BranchName))
	req.Address = strings.TrimSpace(strings.ToUpper(req.Address))
	req.CreatedBy = strings.TrimSpace(req.CreatedBy)
}
func (req *reqCreate) validate() error {
	validationErrors := []errorHandler.HttpError{}
	if !commonvalidator.IsRequired(req.TenantID) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("tenant_id"))
	}
	if !commonvalidator.IsRequired(req.BranchName) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("branch_name"))
	}
	if !commonvalidator.IsRequired(req.Address) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("address"))
	}
	if !commonvalidator.IsRequired(req.CreatedBy) {
		validationErrors = append(validationErrors, *utility.ConstructErrorRequired("created_by"))
	}
	if !commonvalidator.MaxLen(req.BranchName, 40) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("branch_name"))
	}
	if !commonvalidator.MaxLen(req.Address, 500) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("address"))
	}
	if !commonvalidator.MaxLen(req.TenantID, 500) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("tenant_id"))
	}
	if !commonvalidator.MaxLen(req.CreatedBy, 100) {
		validationErrors = append(validationErrors, *utility.ConstructErrorMaxLen("created_by"))
	}
	if len(validationErrors) > 0 {
		return errorHandler.NewMultipleFieldsValidation(validationErrors)
	}
	return nil
}
func (s *Service) Create(ctx context.Context, request *RequestCreate) error {
	input := &reqCreate{request}
	input.sanitize()
	if err := input.validate(); err != nil {
		return err
	}
	isOwner, _, err := commonvalidator.IsAllowedUser(
		ctx, s.userStore, input.CreatedBy)
	if err != nil {
		return err
	}
	_, err = commonvalidator.ValidateTenantOwnership(
		ctx, s.tenantStore, input.TenantID, input.CreatedBy)
	if err != nil {
		return err
	}
	existingBranch, err := s.branchStore.FindOne(ctx,
		&store.BranchesFilter{Name: request.BranchName, TenantID: input.TenantID})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existingBranch != nil {
		return errorHandler.NewBadRequest(errorHandler.WithInfo("branch is exist"))
	}
	totalBranchesPerTenant, err := s.branchStore.TotalBranches(ctx, input.TenantID)
	if err != nil {
		return err
	}
	if totalBranchesPerTenant >= uint8(s.cfg.MaxBranchesPerTenant) {
		errMsg := fmt.Sprintf("max branch per tenant is %d", s.cfg.MaxBranchesPerTenant)
		return errorHandler.NewBadRequest(errorHandler.WithInfo(errMsg))
	}
	insertedNewBranchData := &store.BranchesData{
		BranchID:  uuid.NewString(),
		TenantID:  request.TenantID,
		Name:      request.BranchName,
		Address:   request.Address,
		CreatedAt: time.Now().UTC(),
	}
	if isOwner {
		insertedNewBranchData.CreatedByOwner = sql.NullString{String: request.CreatedBy, Valid: true}
	}
	err = sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return s.branchStore.Insert(txContext, insertedNewBranchData)
	})
	if err != nil {
		return err
	}
	return nil
}
