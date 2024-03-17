package auth

import (
	"context"
	"testing"

	payload "auth-service-rizkysr90-pos/internal/payload/http/auth"
	"auth-service-rizkysr90-pos/internal/store"
	jwttoken "auth-service-rizkysr90-pos/pkg/jwt"

	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {

	ctx := context.Background()
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDB.ExpectBegin()
	mockDB.ExpectCommit()
	mockStore := new(store.UserStoreMock)
	requestPayload := &payload.ReqCreateAccount{
		FirstName:       "Rizki",
		LastName:        "Ramadhan",
		Email:           "notregistered@gmail.com",
		Password:        "verysecretpassword",
		ConfirmPassword: "verysecretpassword",
	}
	queryFilter := &store.UserFilterBy{
		Email: requestPayload.Email,
	}
	// set expectation when user.FindOne is called
	mockStore.On("FindOne", queryFilter, "findactiveuser").Return(nil, nil)
	mockStore.On("Create", mock.Anything).Return(nil)
	jwtTest := jwttoken.New()
	authService := NewAuthService(db, mockStore, jwtTest)
	res := authService.CreateUser(ctx, requestPayload)
	mockStore.AssertExpectations(t)
	assert.Nil(t, res)
}
func TestCreateUserDuplicate(t *testing.T) {
	ctx := context.Background()
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDB.ExpectBegin()
	mockDB.ExpectCommit()
	mockStore := new(store.UserStoreMock)
	requestPayload := &payload.ReqCreateAccount{
		FirstName:       "Rizki",
		LastName:        "Ramadhan",
		Email:           "registered@gmail.com",
		Password:        "verysecretpassword",
		ConfirmPassword: "verysecretpassword",
	}
	queryFilter := &store.UserFilterBy{
		Email: requestPayload.Email,
	}
	// set expectation when user.FindOne is called
	mockStore.On("FindOne", queryFilter, "findactiveuser").Return(&store.UserData{}, nil)
	jwtTest := jwttoken.New()
	authService := NewAuthService(db, mockStore, jwtTest)
	res := authService.CreateUser(ctx, requestPayload)
	mockStore.AssertExpectations(t)
	assert.Error(t, res)
}
func TestCreateUserInvalidPassword(t *testing.T) {
	ctx := context.Background()
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDB.ExpectBegin()
	mockDB.ExpectCommit()
	mockStore := new(store.UserStoreMock)
	requestPayload := &payload.ReqCreateAccount{
		FirstName:       "Rizki",
		LastName:        "Ramadhan",
		Email:           "registered@gmail.com",
		Password:        "verysecretpassword",
		ConfirmPassword: "thedifferent",
	}
	jwtTest := jwttoken.New()
	authService := NewAuthService(db, mockStore, jwtTest)
	res := authService.CreateUser(ctx, requestPayload)
	mockStore.AssertExpectations(t)
	assert.Error(t, res)
}

func TestLoginUser(t *testing.T) {
	ctx := context.Background()
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mockDB.ExpectBegin()
	mockDB.ExpectCommit()
	defer db.Close()
	mockStore := store.UserStoreMock{}
	reqPayload := &payload.ReqLoginUser{
		Email:    "testloginuser@gmail.com",
		Password: "mysecretpassword",
	}
	queryFilter := &store.UserFilterBy{
		Email: reqPayload.Email,
	}
	expectedResQuery := &store.UserData{
		ID:       "89189243-ac3f-46e4-90b5-50ee07300756",
		Email:    "testloginuser@gmail.com",
		Password: "$2a$10$h2rEwjxtJPi2ePYAmvOQr.eDMbcNZxaRrEwLiir9lfFBH9yIBdK3u",
	}
	mockStore.On("FindOne", queryFilter, "findactiveuser").Return(expectedResQuery, nil)
	mockStore.On("Update", mock.Anything)
	mockJWT := new(jwttoken.MockJWTToken)
	authService := NewAuthService(db, &mockStore, mockJWT)
	// var res *payload.ResLoginUser
	mockJWT.On("Generate", &jwttoken.JWTClaims{UserID: expectedResQuery.ID}).Return(
		"jwttoken", nil,
	)
	mockJWT.On("GenerateRefreshToken", &jwttoken.JWTClaims{UserID: expectedResQuery.ID}).Return(
		"jwttoken", nil,
	)
	var res *payload.ResLoginUser
	res, err = authService.LoginUser(ctx, reqPayload)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.RefreshToken)
	mockStore.AssertExpectations(t)
}
func TestLoginUserInvalidPassword(t *testing.T) {
	ctx := context.Background()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockStore := store.UserStoreMock{}
	reqPayload := &payload.ReqLoginUser{
		Email:    "wrongpassword@gmail.com",
		Password: "mysecretpassword",
	}
	mockJWT := new(jwttoken.MockJWTToken)
	svc := NewAuthService(db, &mockStore, mockJWT)
	mockStore.On("FindOne", &store.UserFilterBy{
		Email: reqPayload.Email,
	}, "findactiveuser").Return(&store.UserData{
		Password: "$2a$10$1eo1tHKWjbol1fl512xLNO6B4zI8fsY7Sw1gsa4D8nEzXj1gSkyz2",
	}, nil)
	var res *payload.ResLoginUser
	res, err = svc.LoginUser(ctx, reqPayload)
	assert.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	assert.Nil(t, res)
	mockStore.AssertExpectations(t)
}
func TestLoginUserNotFound(t *testing.T) {
	ctx := context.Background()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockStore := store.UserStoreMock{}
	reqPayload := &payload.ReqLoginUser{
		Email:    "usernotfound@gmail.com",
		Password: "mysecretpassword",
	}
	mockJWT := new(jwttoken.MockJWTToken)
	svc := NewAuthService(db, &mockStore, mockJWT)
	mockStore.On("FindOne", &store.UserFilterBy{
		Email: reqPayload.Email,
	}, "findactiveuser").Return(nil, pgx.ErrNoRows)
	var res *payload.ResLoginUser
	res, err = svc.LoginUser(ctx, reqPayload)
	assert.Contains(t, err.Error(), "user not found")
	assert.Nil(t, res)
	assert.Error(t, err)
	mockStore.AssertExpectations(t)
}
func TestRefreshTokenSucess(t *testing.T) {
	ctx := context.Background()
	const userID = "50c8d653-4a6a-45cf-92fa-406492b463d7"
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockStore := store.UserStoreMock{}
	mockStore.On("FindOne", mock.Anything, "findRefreshToken").Return(&store.UserData{
		ID: userID,
	}, nil)
	reqPayload := &payload.ReqRefreshToken{
		RefreshToken: "refreshtoken",
	}
	mockJWT := new(jwttoken.MockJWTToken)
	mockJWT.On("AuthorizeRefreshToken", reqPayload.RefreshToken).Return(&jwttoken.JWTClaims{
		UserID: userID,
	}, nil)
	mockJWT.On("Generate", &jwttoken.JWTClaims{
		UserID: userID,
	}).Return("access_token", nil)
	svc := NewAuthService(db, &mockStore, mockJWT)
	var res *payload.ResRefreshToken
	res, err = svc.RefreshToken(ctx, reqPayload)
	mockJWT.AssertExpectations(t)
	mockStore.AssertExpectations(t)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.AccessToken)
}

func TestRefreshTokenInvalid(t *testing.T) {
	ctx := context.Background()
	const userID = "50c8d653-4a6a-45cf-92fa-406492b463d7"
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockStore := store.UserStoreMock{}
	mockStore.On("FindOne", mock.Anything, "findRefreshToken").Return(&store.UserData{
		ID: userID,
	}, nil)
	reqPayload := &payload.ReqRefreshToken{
		RefreshToken: "stolenToken",
	}
	mockJWT := new(jwttoken.MockJWTToken)
	mockJWT.On("AuthorizeRefreshToken", reqPayload.RefreshToken).Return(&jwttoken.JWTClaims{
		UserID: "50c8d653-4a6a-45cf-92fa-406492b463d8", // DIFF ID
	}, nil)
	svc := NewAuthService(db, &mockStore, mockJWT)
	var res *payload.ResRefreshToken
	res, err = svc.RefreshToken(ctx, reqPayload)
	mockJWT.AssertExpectations(t)
	mockStore.AssertExpectations(t)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	expectedError := restapierror.NewUnauthorized(restapierror.WithMessage("invalid token"))
	assert.Equal(t, expectedError.Error(), err.Error())
}
