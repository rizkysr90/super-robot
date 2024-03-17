package jwttoken

type JWTInterface interface {
	Generate(jwtClaims *JWTClaims) (string, error)
	GenerateRefreshToken(jwtClaims *JWTClaims) (string, error)
	Authorize(tokenString string) error
}
