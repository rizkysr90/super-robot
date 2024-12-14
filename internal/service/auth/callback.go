package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rizkysr90-pos/internal/constant"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/internal/utility"
	"rizkysr90-pos/pkg/errorHandler"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/rizkysr90/rizkysr90-go-pkg/sqldb"
	"golang.org/x/oauth2"
)

type RequestCallback struct {
	State string `validate:"required"`
	Code  string `validate:"required"`
}
type ResponseCallback struct {
	SessionData *store.SessionRedisData
}
type requestCallback struct {
	payload *RequestCallback
	auth    *Auth
}
type UserInfoClaims struct {
	StateData *store.StateData
	Email     string `json:"email"`
	Name      string `json:"name"`
	Sub       string `json:"sub"`
}

func (req *requestCallback) sanitize() {
	req.payload.State = strings.TrimSpace(req.payload.State)
	req.payload.Code = strings.TrimSpace(req.payload.Code)
}
func (req *requestCallback) validate() error {
	validationUtil := utility.NewValidationUtil()
	return validationUtil.Validate(req.payload)
}
func (req *requestCallback) getAndVerifyState(ctx context.Context) (*store.StateData, error) {
	stateData, err := req.auth.stateStore.FindOne(ctx, req.payload.State)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("stateID is invalid or not found")
		}
		return nil, fmt.Errorf("failed to get state data, got : %w", err)
	}
	return stateData, nil
}
func (req *requestCallback) exchangeToken(ctx context.Context) (*oauth2.Token, error) {
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("grant_type", "authorization_code"),
	}
	var oauth2Token *oauth2.Token
	oauth2Token, err := req.auth.authClient.Oauth.Exchange(ctx, req.payload.Code, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange the token, got : %w", err)
	}
	return oauth2Token, nil
}
func (req *requestCallback) getAndVerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	// Get raw ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token in response")
	}
	idToken, err := req.auth.authClient.OIDC.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id token, got : %w", err)
	}
	return idToken, nil
}
func (req *requestCallback) getUserInfoData(idToken *oidc.IDToken) (*UserInfoClaims, error) {
	UserInfoClaims := &UserInfoClaims{}
	if err := idToken.Claims(&UserInfoClaims); err != nil {
		return nil, fmt.Errorf("failed to extract claims, got : %w", err)
	}
	return UserInfoClaims, nil
}
func (req *requestCallback) checkExistingUser(ctx context.Context, email string) (*store.UserData, error) {
	userQueryFilter := &store.UserQueryFilter{Email: email}
	checkUserData, err := req.auth.userStore.FindOne(ctx, userQueryFilter)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return checkUserData, nil
}
func (req *requestCallback) setSessionData(userInfoData *UserInfoClaims) (*store.SessionRedisData, error) {
	setUserID := uuid.NewString()
	setTenantID := uuid.NewString()
	generateSessionID, err := utility.GenerateRandomBase64Str()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session id, got : %w", err)
	}
	sessionData := &store.SessionRedisData{
		SessionID:    generateSessionID,
		UserID:       setUserID,
		UserType:     "headoffice",
		UserEmail:    userInfoData.Email,
		UserFullName: userInfoData.Name,
		UserTenantID: setTenantID,
		UserBranchID: "",
		UserAuthType: string(constant.UserTypeGoogle),
		UserRoles:    []string{},
		CreatedAt:    time.Now().UTC(),
	}
	return sessionData, nil
}
func (req *requestCallback) onlyUpdateUser(txContext context.Context, existingUser *store.UserData) error {
	updatedData := &store.UserData{
		ID:          existingUser.ID,
		LastLoginAt: time.Now().UTC(),
		Email:       existingUser.Email,
		FullName:    existingUser.FullName,
		GoogleID:    existingUser.GoogleID,
	}
	if err := req.auth.userStore.Update(txContext, updatedData); err != nil {
		return fmt.Errorf("failed to update user, got: %w", err)
	}
	return nil
}
func (req *requestCallback) createNewUser(txContext context.Context,
	sessionData *store.SessionRedisData, userInfoData *UserInfoClaims) error {
	insertedTenantData := &store.TenantData{
		ID:        sessionData.UserTenantID,
		Name:      userInfoData.StateData.TenantName.String,
		OwnerID:   sql.NullString{String: "", Valid: false},
		CreatedAt: time.Now().UTC(),
	}
	insertedUserData := &store.UserData{
		ID:           sessionData.UserID,
		Email:        userInfoData.Email,
		FullName:     userInfoData.Name,
		GoogleID:     sql.NullString{String: userInfoData.Sub, Valid: true},
		PasswordHash: sql.NullString{String: "", Valid: false},
		AuthType:     "google",
		UserType:     "owner",
		TenantID:     sessionData.UserTenantID,
		CreatedAt:    time.Now().UTC(),
		LastLoginAt:  time.Now().UTC(),
	}
	updatedTenantData := &store.TenantData{
		ID:        sessionData.UserTenantID,
		OwnerID:   sql.NullString{String: sessionData.UserID, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	}
	if err := req.auth.tenantStore.Insert(txContext, insertedTenantData); err != nil {
		return fmt.Errorf("failed to insert tenant data, got : %w", err)
	}
	if err := req.auth.userStore.Insert(txContext, insertedUserData); err != nil {
		return fmt.Errorf("failed to insert user data, got : %w", err)
	}
	if err := req.auth.tenantStore.Update(txContext, updatedTenantData); err != nil {
		return fmt.Errorf("failed to update tenant data, got : %w", err)
	}
	return nil
}
func (a *Auth) Callback(ctx context.Context, request *RequestCallback) (*ResponseCallback, error) {
	input := &requestCallback{payload: request, auth: a}
	input.sanitize()
	if err := input.validate(); err != nil {
		return nil, err
	}
	userInfoData, err := input.handleCallback(ctx)
	if err != nil {
		return nil, err
	}
	checkUserData, err := input.checkExistingUser(ctx, userInfoData.Email)
	if err != nil {
		return nil, err
	}
	isSameEmail := checkUserData != nil
	sessionData, err := input.setSessionData(userInfoData)
	if err != nil {
		return nil, err
	}
	if err = a.session.Insert(ctx, sessionData); err != nil {
		return nil, fmt.Errorf("failed to set redis data, got : %w", err)
	}
	err = sqldb.WithinTx(ctx, a.db, func(tx sqldb.QueryExecutor) error {
		txContext := sqldb.WithTxContext(ctx, tx)
		if err = a.stateStore.Delete(txContext, userInfoData.StateData.ID); err != nil {
			return fmt.Errorf("failed to delete state data, got : %w", err)
		}
		if isSameEmail {
			return input.onlyUpdateUser(txContext, checkUserData)
		}
		return input.createNewUser(txContext, sessionData, userInfoData)
	})
	return &ResponseCallback{SessionData: sessionData}, nil
}
func (req *requestCallback) handleCallback(ctx context.Context) (*UserInfoClaims, error) {
	stateData, err := req.getAndVerifyState(ctx)
	if err != nil {
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo(fmt.Sprintf("failed to get state, got : %s", err.Error())))
	}
	oauth2Token, err := req.exchangeToken(ctx)
	if err != nil {
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo(fmt.Sprintf("failed to exchange token, got : %s", err.Error())))
	}
	idToken, err := req.getAndVerifyIDToken(ctx, oauth2Token)
	if err != nil {
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo(fmt.Sprintf("failed to get or verify id token, got : %s", err.Error())))
	}
	userInfoData, err := req.getUserInfoData(idToken)
	if err != nil {
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo(fmt.Sprintf("failed to get user info data, got : %s", err.Error())))
	}
	userInfoData.StateData = stateData
	return userInfoData, nil
}
