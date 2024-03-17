package jwttoken

import "github.com/stretchr/testify/mock"

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
