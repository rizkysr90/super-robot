package store

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type UserStoreMock struct {
	mock.Mock
}

func (u *UserStoreMock) Create(ctx context.Context, createdData *User) error {
	u.Called(createdData)
	return nil
}
