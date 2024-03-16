package store

import (
	"context"
	"fmt"

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
	args := u.Called(filterBy, staging)
	if filterBy.Email == "notregistered@gmail.com" {
		return nil, nil
	}
	if filterBy.Email == "testloginuser@gmail.com" {
		// Assuming args.Get(0) returns an interface{}
		value := args.Get(0)

		// Perform type assertion
		if mergeData, ok := value.(*UserData); ok {
			// Type assertion successful, mergeData is now of type *MergeReportsData
			return mergeData, args.Error(1)
		} else {
			// Type assertion failed
			return nil, fmt.Errorf("unexpected type, expected *MergeReportsData")
		}
	}
	if filterBy.Email == "wrongpassword@gmail.com" {
		// Assuming args.Get(0) returns an interface{}
		value := args.Get(0)

		// Perform type assertion
		if _, ok := value.(*UserData); ok {
			// Type assertion successful, userData is now of type *UserData
			return &UserData{
					Password: "$2a$10$1eo1tHKWjbol1fl512xLNO6B4zI8fsY7Sw1gsa4D8nEzXj1gSkyz2"},
				args.Error(1)
		} else {
			// Type assertion failed
			return nil, fmt.Errorf("unexpected type, expected *UserData")
		}
	}
	if filterBy.Email == "usernotfound@gmail.com" {
		return nil, args.Error(1)
	}
	return &UserData{}, nil
}
