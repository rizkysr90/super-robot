package users

import (
	payload "auth-service-rizkysr90-pos/internal/payload/http/users"
	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
)

func (u *UserService) LoginUser(ctx context.Context, 
	req *payload.ReqLoginUsers) (*payload.ResLoginUsers, error) {
	
	user, err := u.userStore.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, restapierror.NewNotFound(restapierror.WithMessage("user not found"))
		}
		return nil, err
	}
	// Compare password with bcrypt
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return nil, restapierror.NewBadRequest()
    }
	accessToken, err := u.jwt.Generate(&jwttoken.JWTClaims{UserID: user.ID})
	if err != nil {
		return nil, err
	}
	err = sqldb.WithinTx(ctx, u.db, func(qe sqldb.QueryExecutor) error {
		tx := sqldb.WithTxContext(ctx, qe)
		return u.userStore.UpdateUserAccessToken(tx, req.Username, accessToken )

	})
	if err != nil {
		return nil, err
	}
	return &payload.ResLoginUsers{AccessToken: accessToken}, nil
}