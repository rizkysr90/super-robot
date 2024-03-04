package auth

import (
	"context"
	"testing"

	payload "github.com/rizkysr90/go-boilerplate/internal/payload/http/auth"
	"github.com/rizkysr90/go-boilerplate/internal/store"

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

	authService := NewAuthService(db, mockStore)
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

	authService := NewAuthService(db, mockStore)
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
	authService := NewAuthService(db, mockStore)
	res := authService.CreateUser(ctx, requestPayload)
	mockStore.AssertExpectations(t)
	assert.Error(t, res)
}
