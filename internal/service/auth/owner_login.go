package auth

import (
	"context"
	"fmt"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/utility"
	"rizkysr90-pos/pkg/errorHandler"

	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

func (a *Auth) LoginOwner(ctx context.Context) (string, error) {
	stateID, err := utility.GenerateRandomBase64Str()
	if err != nil {
		return "", errorHandler.NewInternalServer(errorHandler.WithInfo(
			fmt.Sprintf("failed to generate state id, got : %s", err.Error()),
		))
	}
	insertedStateData := store.StateData{ID: stateID}
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
