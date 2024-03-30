package store

import (
	payload "auth-service-rizkysr90-pos/internal/payload/http/store"
	"auth-service-rizkysr90-pos/internal/store"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

type Service struct {
	db         *sql.DB
	storeStore store.StoreStore
}

func NewStoreService(db *sql.DB, storeStore store.StoreStore) *Service {
	return &Service{
		db:         db,
		storeStore: storeStore,
	}
}
func (s *Service) CreateStore(ctx context.Context, req *payload.ReqCreateStore) error {
	insertedData := store.StoreData{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Address:   req.Address,
		Contact:   req.Contact,
		UserID:    req.UserID,
		CreatedAt: time.Now().UTC(),
	}
	err := sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return s.storeStore.Insert(txContext, &insertedData)
	})
	if err != nil {
		return err
	}

	return nil
}
