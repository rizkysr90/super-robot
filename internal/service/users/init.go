package users

import (
	"database/sql"

	"auth-service-rizkysr90-pos/internal/store"
	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"
)

type Service struct {
	db        *sql.DB
	userStore store.UserStore
	jwtToken  jwttoken.JWTInterface
}

func NewUsersService(db *sql.DB, userStore store.UserStore,
	jwttoken jwttoken.JWTInterface) *Service {
	return &Service{
		db:        db,
		userStore: userStore,
		jwtToken:  jwttoken,
	}
}

// func (s *Service) LoginUser(ctx context.Context,
// 	req *payload.ReqLoginUser) (*payload.ResLoginUser, error) {
// 	result, err := s.userStore.FindOne(ctx, &store.UserFilterBy{
// 		Email: req.Email,
// 	}, "findactiveuser")
// 	if err != nil {
// 		if strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
// 			return nil, restapierror.NewNotFound(restapierror.WithMessage("user not found"))
// 		}
// 	}
// 	if err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(req.Password)); err != nil {
// 		if err.Error() == bcrypt.ErrMismatchedHashAndPassword.Error() {
// 			return nil, restapierror.NewBadRequest(restapierror.WithMessage("invalid username or password"))
// 		}
// 		return nil, restapierror.NewBadRequest(restapierror.WithMessage(err.Error()))
// 	}
// 	// gen token
// 	var genToken string
// 	genToken, err = s.jwtToken.Generate(&jwttoken.JWTClaims{
// 		UserID: result.ID,
// 		Role:   1,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	// gen refresh token
// 	var refreshToken string
// 	refreshToken, err = s.jwtToken.GenerateRefreshToken(&jwttoken.JWTClaims{
// 		UserID: result.ID,
// 		Role:   1,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = sqldb.WithinTx(ctx, s.db, func(tx sqldb.QueryExecutor) error {
// 		txContext := sqldb.WithTxContext(ctx, tx)
// 		return s.userStore.Update(txContext, &store.UserData{
// 			RefreshToken: refreshToken,
// 		}, &store.UserFilterBy{Email: req.Email}, "updaterefreshtoken")

// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &payload.ResLoginUser{
// 		Token:        genToken,
// 		RefreshToken: refreshToken,
// 	}, nil
// }

// func (s *Service) RefreshToken(ctx context.Context,
// 	req *payload.ReqRefreshToken) (*payload.ResRefreshToken, error) {
// 	claims, err := s.jwtToken.AuthorizeRefreshToken(req.RefreshToken)
// 	if err != nil {
// 		return nil, restapierror.NewUnauthorized(restapierror.WithMessage("invalid token"))
// 	}
// 	// find by token
// 	var data *store.UserData
// 	data, err = s.userStore.FindOne(ctx, &store.UserFilterBy{
// 		RefreshToken: req.RefreshToken,
// 	}, "findRefreshToken")
// 	if err != nil {
// 		if sql.ErrNoRows.Error() == err.Error() {
// 			return nil, restapierror.NewUnauthorized(restapierror.WithMessage("invalid token"))
// 		}
// 		return nil, err
// 	}
// 	if data.ID != claims.Subject {
// 		return nil, restapierror.NewUnauthorized(restapierror.WithMessage("invalid token"))
// 	}
// 	// gen access token
// 	var accessToken string
// 	accessToken, err = s.jwtToken.Generate(&jwttoken.JWTClaims{
// 		UserID: data.ID,
// 		Role:   1,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &payload.ResRefreshToken{AccessToken: accessToken}, nil
// }
