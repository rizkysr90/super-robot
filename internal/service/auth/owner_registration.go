package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/utility"
	"rizkysr90-pos/pkg/errorHandler"
	"strings"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type RequestRegisterOwner struct {
	TenantName string `json:"tenant_name"`
}
type requestRegisterOwner struct {
	payload *RequestRegisterOwner
}

func (req *requestRegisterOwner) sanitize() {
	req.payload.TenantName = strings.TrimSpace(req.payload.TenantName)
}
func (req *requestRegisterOwner) validate() error {
	if len(req.payload.TenantName) > 100 {
		return errors.New("max tenant name length is 100 characters")
	}
	return nil
}
func (a *Auth) RegisterOwner(ctx context.Context, request *RequestRegisterOwner) (string, error) {
	input := &requestRegisterOwner{payload: request}
	input.sanitize()
	if err := input.validate(); err != nil {
		return "", errorHandler.NewBadRequest(errorHandler.WithInfo(
			fmt.Sprintf("failed to validate request, got : %s", err.Error()),
		))
	}
	stateID, err := utility.GenerateRandomBase64Str()
	if err != nil {
		return "", errorHandler.NewInternalServer(errorHandler.WithInfo(
			fmt.Sprintf("failed to generate state id, got : %s", err.Error()),
		))
	}
	insertedStateData := store.StateData{
		ID:         stateID,
		TenantName: sql.NullString{String: input.payload.TenantName, Valid: true},
	}
	err = sqldb.WithinTx(ctx, a.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return a.stateStore.Insert(txContext, &insertedStateData)
	})
	if err != nil {
		return "", errorHandler.NewInternalServer(errorHandler.WithInfo(
			fmt.Sprintf("failed to insert state data, got : %s", err.Error()),
		))
	}
	return stateID, nil
}
