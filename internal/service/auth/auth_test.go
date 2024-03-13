package auth

import (
	"context"
	"fmt"
	"os"
	"testing"

	payload "github.com/rizkysr90/go-boilerplate/internal/payload/http/auth"
	"github.com/rizkysr90/go-boilerplate/internal/store"
	jwttoken "github.com/rizkysr90/go-boilerplate/pkg/jwt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var privateKeyTest string
var publicKeyTest string

func TestMain(m *testing.M) {
	var err error
	var privateKey []byte
	var publicKey []byte
	privateKey, err = os.ReadFile("private_key.pem")
	if err != nil {
		fmt.Println("Error loading private key:", err)
	}
	privateKeyTest = string(privateKey)
	publicKey, err = os.ReadFile("public_key.pem")
	if err != nil {
		fmt.Println("Error loading public key:", err)
	}
	publicKeyTest = string(publicKey)
	m.Run()
}
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
