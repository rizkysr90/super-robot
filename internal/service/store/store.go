package store

import (
	payload "auth-service-rizkysr90-pos/internal/payload/http/store"
	"auth-service-rizkysr90-pos/internal/store"
	"context"
	"database/sql"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"
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

func (s *Service) CreateStore(ctx context.Context, req *payload.ReqCreateStore) (*payload.ResCreateStore, error) {
	// check duplicate store name per user
	_, err := s.storeStore.Finder(ctx, &store.StoreFilter{Name: req.Name, EmployeeID: req.EmployeeID}, "findbyname")
	if err != nil && err != sql.ErrNoRows {
		if !strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return nil, restapierror.NewInternalServer(restapierror.WithMessage(err.Error()))
		}
	}
	// if error its not nil, it means there aru duplicate name
	if err == nil {
		return nil, restapierror.NewBadRequest(restapierror.WithMessage("duplicate store name"))
	}
	insertedData := store.StoreData{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Address:   req.Address,
		Contact:   req.Contact,
		UserID:    req.EmployeeID,
		CreatedAt: time.Now().UTC(),
	}
	err = sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return s.storeStore.Insert(txContext, &insertedData)
	})
	if err != nil {
		return nil, err
	}

	return &payload.ResCreateStore{StoreID: insertedData.ID, CreatedAt: insertedData.CreatedAt.String()}, nil
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
		EmployeeID: req.EmployeeID,
		Pagination: store.Pagination{
			PageSize:   req.PageSize,
			PageNumber: offset,
		},
	}
	data, err := s.storeStore.Finder(ctx, queryFilter, "getallstore")
	if err != nil {
		return nil, err
	}
	var setResponseData []store.SetResponse
	// Type assertion to retrieve the underlying data value
	getData, ok := data.([]store.StoreData)
	if ok {
		for _, val := range getData {
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
	} else {
		return nil, restapierror.NewInternalServer(restapierror.WithMessage("invalid type assertion"))
	}
}
func (s *Service) DeleteStore(ctx context.Context, req *payload.ReqDeleteStore) error {
	setFilter := store.StoreFilter{
		StoreID:    req.StoreID,
		EmployeeID: req.EmployeeID,
	}
	err := sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return s.storeStore.Delete(txContext, &setFilter)
	})
	if err != nil {
		return restapierror.NewInternalServer(restapierror.WithMessage(err.Error()))
	}
	return nil
}
func (s *Service) UpdateStore(ctx context.Context, req *payload.ReqUpdateStore) error {
	// check duplicate store name per user
	_, err := s.storeStore.Finder(ctx, &store.StoreFilter{Name: req.Name, EmployeeID: req.EmployeeID}, "findbyname")
	if err != nil && err != sql.ErrNoRows {
		if !strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return restapierror.NewInternalServer(restapierror.WithMessage(err.Error()))
		}
	}

	// if error its not nil, it means there aru duplicate name
	if err == nil {
		return restapierror.NewBadRequest(restapierror.WithMessage("duplicate store name"))
	}
	updatedData := store.StoreData{
		Name:    req.Name,
		Address: req.Address,
		Contact: req.Contact,
	}
	filterData := store.StoreFilter{
		StoreID:    req.StoreID,
		EmployeeID: req.EmployeeID,
	}
	err = sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		var result int64
		result, err = s.storeStore.Update(txContext, &updatedData, &filterData, "updatestore")
		if result == 0 {
			return restapierror.NewNotFound(restapierror.WithMessage("data not found"))
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
