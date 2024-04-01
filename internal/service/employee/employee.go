package employee

import (
	payload "auth-service-rizkysr90-pos/internal/payload/http/employee"
	"auth-service-rizkysr90-pos/internal/store"
	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db            *sql.DB
	employeeStore store.EmployeeStore
	jwtToken      jwttoken.JWTInterface
}

func NewEmployeeService(db *sql.DB, employeeStore store.EmployeeStore,
	jwtTOken jwttoken.JWTInterface) *Service {
	return &Service{
		db:            db,
		employeeStore: employeeStore,
		jwtToken:      jwtTOken,
	}
}
func (s *Service) Create(ctx context.Context, req *payload.ReqCreateEmployee) error {
	// gen bcrypt
	bytesPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return err
	}
	insertedData := store.EmployeeData{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Contact:   req.Contact,
		Username:  req.Username,
		Password:  string(bytesPassword),
		StoreID:   req.StoreID,
		Role:      req.Role,
		CreatedAt: time.Now().UTC(),
	}
	err = sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return s.employeeStore.Insert(txContext, &insertedData)
	})
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) Login(ctx context.Context,
	req *payload.ReqLoginEmployee) (*payload.ResLoginEmployee, error) {
	result, err := s.employeeStore.FindOne(ctx, &store.EmployeeFilter{
		Username: req.Username,
	}, "findactiveemployee")
	if err != nil {
		if strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
			return nil, restapierror.NewNotFound(restapierror.WithMessage("user not found"))
		}
	}
	if err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(req.Password)); err != nil {
		if err.Error() == bcrypt.ErrMismatchedHashAndPassword.Error() {
			return nil, restapierror.NewBadRequest(restapierror.WithMessage("invalid username or password"))
		}
		return nil, restapierror.NewBadRequest(restapierror.WithMessage(err.Error()))
	}
	// gen token
	var genToken string
	jwtClaims := &jwttoken.JWTClaims{
		UserID: result.ID,
		Role:   result.Role,
	}
	genToken, err = s.jwtToken.Generate(jwtClaims)
	if err != nil {
		return nil, err
	}
	// gen refresh token
	var refreshToken string
	refreshToken, err = s.jwtToken.GenerateRefreshToken(jwtClaims)
	if err != nil {
		return nil, err
	}
	err = sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		return s.employeeStore.Update(txContext, &store.EmployeeData{
			RefreshToken: refreshToken,
		}, &store.EmployeeFilter{Username: req.Username}, "updaterefreshtoken")

	})
	if err != nil {
		return nil, err
	}
	return &payload.ResLoginEmployee{
		Token:        genToken,
		RefreshToken: refreshToken,
		Role:         result.Role,
	}, nil
}
func (s *Service) RefreshToken(ctx context.Context,
	req *payload.ReqRefreshToken) (*payload.ResRefreshToken, error) {
	claims, err := s.jwtToken.AuthorizeRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, restapierror.NewUnauthorized(restapierror.WithMessage("invalid token"))
	}
	// find by token
	var data *store.EmployeeData
	data, err = s.employeeStore.FindOne(ctx, &store.EmployeeFilter{
		ID: claims.Subject,
	}, "findByID")
	if err != nil {
		if sql.ErrNoRows.Error() == err.Error() {
			return nil, restapierror.NewUnauthorized(restapierror.WithMessage("invalid token"))
		}
		return nil, err
	}
	if data.RefreshToken != req.RefreshToken {
		return nil, restapierror.NewUnauthorized(restapierror.WithMessage("invalid token"))
	}
	// gen access token
	var accessToken string
	accessToken, err = s.jwtToken.Generate(&jwttoken.JWTClaims{
		UserID: data.ID,
		Role:   data.Role,
	})
	if err != nil {
		return nil, err
	}
	return &payload.ResRefreshToken{AccessToken: accessToken, Role: data.Role}, nil
}

// func (s *Service) GetAllStore(ctx context.Context, req *payload.ReqGetAllStore) (*payload.ResGetAllStore, error) {
// 	if req.PageSize == 0 {
// 		req.PageSize = 50
// 	}
// 	if req.PageNumber == 0 {
// 		req.PageNumber = 1
// 	}
// 	offset := (req.PageNumber - 1) * req.PageSize
// 	queryFilter := &store.StoreFilter{
// 		UserID: req.UserID,
// 		Pagination: store.Pagination{
// 			PageSize:   req.PageSize,
// 			PageNumber: offset,
// 		},
// 	}
// 	data, err := s.storeStore.Finder(ctx, queryFilter, "getallpagination")
// 	if err != nil {
// 		return nil, err
// 	}
// 	var setResponseData []store.SetResponse
// 	for _, val := range data {
// 		setResponse := val.ToPayloadResponse(&val)
// 		setResponseData = append(setResponseData, *setResponse)
// 	}
// 	response := &payload.ResGetAllStore{
// 		Data: setResponseData,
// 		Pagination: store.Pagination{
// 			PageSize:   req.PageSize,
// 			PageNumber: req.PageNumber,
// 			TotalPages: int(math.Ceil(
// 				float64(queryFilter.Pagination.TotalElements) /
// 					float64(queryFilter.Pagination.PageSize))),
// 			TotalElements: queryFilter.Pagination.TotalElements,
// 		},
// 	}
// 	return response, nil
// }
