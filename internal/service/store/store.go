package store

import (
	payload "auth-service-rizkysr90-pos/internal/payload/http/store"
	"auth-service-rizkysr90-pos/internal/store"
	"context"
	"database/sql"
	"math"
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
func (s *Service) GetAllStore(ctx context.Context, req *payload.ReqGetAllStore) (*payload.ResGetAllStore, error) {
	if req.PageSize == 0 {
		req.PageSize = 50
	}
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	offset := (req.PageNumber - 1) * req.PageSize
	queryFilter := &store.StoreFilter{
		UserID: req.UserID,
		Pagination: store.Pagination{
			PageSize:   req.PageSize,
			PageNumber: offset,
		},
	}
	data, err := s.storeStore.Finder(ctx, queryFilter, "getallpagination")
	if err != nil {
		return nil, err
	}
	var setResponseData []store.SetResponse
	for _, val := range data {
		setResponse := val.ToPayloadResponse(&val)
		setResponseData = append(setResponseData, *setResponse)
	}
	response := &payload.ResGetAllStore{
		Data: setResponseData,
		Pagination: store.Pagination{
			PageSize:   req.PageSize,
			PageNumber: req.PageNumber,
			TotalPages: int(math.Ceil(
				float64(queryFilter.Pagination.TotalElements) /
					float64(queryFilter.Pagination.PageSize))),
			TotalElements: queryFilter.Pagination.TotalElements,
		},
	}
	return response, nil
}
