package jwttoken

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rizkysr90/rizkysr90-go-pkg/restapierror"
)

type JWT struct {
}

func New() *JWT {
	return &JWT{}
}

type JWTClaims struct {
	UserID string
	Role   int
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
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":  "rizkysr90-pos",
		"sub":  jwtClaims.UserID,
		"exp":  time.Now().Add((time.Hour * 24) * 7).Unix(), // 1 weeks expiry
		"role": jwtClaims.Role,
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
	var privateKeyBytes []byte
	var err error
	privateKeyBytes, err = os.ReadFile("private_key_jwt.pem")
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":  "rizkysr90-pos",
		"sub":  jwtClaims.UserID,
		"exp":  time.Now().Add(time.Minute * 5).Unix(), // 5 minutes expiry
		"role": jwtClaims.Role,
		// "exp": time.Now().Add(time.Second * 2).Unix(), // 2 second

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
func (j *JWT) Authorize(tokenString string) (*MyCustomClaims, error) {
	var err error
	var publicKeyBytes []byte
	var publicKey *rsa.PublicKey
	var jwtToken *jwt.Token
	publicKeyBytes, err = os.ReadFile("public_key_jwt.pem")
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
		return nil, restapierror.NewUnauthorized(restapierror.WithMessage(err.Error()))
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
		return nil, restapierror.NewUnauthorized(restapierror.WithMessage(err.Error()))
	} else if claims, ok := jwtToken.Claims.(*MyCustomClaims); ok {
		return claims, nil
	} else {
		return nil, err
	}
}
