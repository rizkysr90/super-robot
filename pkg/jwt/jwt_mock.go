package jwttoken

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type MockJWTToken struct {
	mock.Mock
}

func (m *MockJWTToken) Generate(jwtClaims *JWTClaims) (string, error) {
	args := m.Called(jwtClaims)
	return args.String(0), args.Error(1)
}
func (m *MockJWTToken) GenerateRefreshToken(jwtClaims *JWTClaims) (string, error) {
	args := m.Called(jwtClaims)
	return args.String(0), args.Error(1)
}
func (m *MockJWTToken) Authorize(tokenString string) error {
	args := m.Called(tokenString)
	return args.Error(0)
}
func (m *MockJWTToken) AuthorizeRefreshToken(tokenString string) (*MyCustomClaims, error) {
	args := m.Called(tokenString)
	if tokenString == "refreshtoken" {
		return &MyCustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "50c8d653-4a6a-45cf-92fa-406492b463d7",
			},
		}, args.Error(1)
	}
	if tokenString == "stolenToken" {
		return &MyCustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "50c8d653-4a6a-45cf-92fa-406492b463d8", // diff userID
			},
		}, args.Error(1)
	}
	return nil, nil
}
