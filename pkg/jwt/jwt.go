package jwttoken

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"rizkysr90-pos/pkg/errorHandler"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secretKey string
}

func New(secretKey string) *JWT {
	return &JWT{secretKey: secretKey}
}

type JWTClaims struct {
	UserID string
}
type MyCustomClaims struct {
	jwt.RegisteredClaims
	JWTClaims
}

func (j *JWT) GenerateRefreshToken(jwtClaims *JWTClaims) (string, error) {
	var privateKeyBytes []byte
	var err error
	privateKeyBytes, err = os.ReadFile("private_key_refresh_jwt.pem")
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "rizkysr90-pos",
		"sub": jwtClaims.UserID,
		"exp": time.Now().Add((time.Hour * 24) * 7).Unix(), // 1 weeks expiry
		// "exp": time.Now().Add(time.Minute * 2).Unix(), // 2 minute expiry

	})
	var signedToken string
	var privateKey *rsa.PrivateKey
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", err
	}
	signedToken, err = token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
func (j *JWT) Generate(jwtClaims *JWTClaims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "rizkysr90-pos",
		"sub": jwtClaims.UserID,
		"exp": time.Now().Add(time.Minute * 5).Unix(), // 5 minutes expiry

	})
	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
func (j *JWT) Authorize(tokenString string) (*MyCustomClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, errorHandler.NewUnauthorized(errorHandler.WithMessage(err.Error()))
	} else if claims, ok := jwtToken.Claims.(*MyCustomClaims); ok {
		return claims, nil
	} else {
		return nil, err
	}
}
func (j *JWT) AuthorizeRefreshToken(tokenString string) (*MyCustomClaims, error) {
	var err error
	var publicKeyBytes []byte
	var publicKey *rsa.PublicKey
	var jwtToken *jwt.Token
	publicKeyBytes, err = os.ReadFile("public_key_refresh_jwt.pem")
	if err != nil {
		return nil, err
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	jwtToken, err = jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, errorHandler.NewUnauthorized(errorHandler.WithMessage(err.Error()))
	} else if claims, ok := jwtToken.Claims.(*MyCustomClaims); ok {
		return claims, nil
	} else {
		return nil, err
	}
}
