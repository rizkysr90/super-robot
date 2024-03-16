package auth

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	payload "github.com/rizkysr90/go-boilerplate/internal/payload/http/auth"
	"github.com/rizkysr90/go-boilerplate/internal/store"
	jwttoken "github.com/rizkysr90/go-boilerplate/pkg/jwt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// var privateKeyTest string
// var publicKeyTest string

//	func TestMain(m *testing.M) {
//		var err error
//		var privateKey []byte
//		var publicKey []byte
//		// Get the current working directory
//		dir, err := os.Getwd()
//		if err != nil {
//			fmt.Println("Error:", err)
//			return
//		}
//		// Construct the file path
//		filePathPriv := filepath.Join(dir, "private_key_test.pem")
//		filePathPub := filepath.Join(dir, "public_key_test.pem")
//		privateKey, err = os.ReadFile(filePathPriv)
//		if err != nil {
//			fmt.Println("Error loading private key:", err)
//		}
//		log.Println(filePathPriv)
//		log.Println(filePathPub)
//		privateKeyTest = string(privateKey)
//		publicKey, err = os.ReadFile(filePathPub)
//		if err != nil {
//			fmt.Println("Error loading public key:", err)
//		}
//		publicKeyTest = string(publicKey)
//		m.Run()
//	}
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
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
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
	mockJWT := new(jwttoken.MockJWTToken)
	authService := NewAuthService(db, &mockStore, mockJWT)
	// var res *payload.ResLoginUser
	mockJWT.On("Generate", &jwttoken.JWTClaims{UserID: expectedResQuery.ID}).Return(
		"jwttoken", nil,
	)
	_, err = authService.LoginUser(ctx, reqPayload)
	assert.Nil(t, err)
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
