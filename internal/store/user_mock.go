package store

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type UserStoreMock struct {
	mock.Mock
}

func (u *UserStoreMock) Create(ctx context.Context, createdData *InsertedData) error {
	u.Called(createdData)
	return nil
}
func (u *UserStoreMock) FindOne(ctx context.Context,
	filterBy *UserFilterBy, staging string) (*UserData, error) {
	u.Called(filterBy, staging)
	if filterBy.Email == "notregistered@gmail.com" {
		return nil, nil
	}
	return &UserData{}, nil
}
