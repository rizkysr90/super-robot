package jwttoken

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
}

func New() *JWT {
	return &JWT{}
}

type JWTClaims struct {
	UserID string
}
type MyCustomClaims struct {
	jwt.RegisteredClaims
}

func (j *JWT) Generate(jwtClaims *JWTClaims) (string, error) {
	var privateKeyBytes []byte
	var err error
	privateKeyBytes, err = os.ReadFile("private_key_jwt.pem")
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "rizkysr90-pos",
		"sub": jwtClaims.UserID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 1 day expiry token
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

func (j *JWT) Authorize(tokenString string) error {
	var err error
	var publicKeyBytes []byte
	var publicKey *rsa.PublicKey
	var jwtToken *jwt.Token
	publicKeyBytes, err = os.ReadFile("public_key_jwt.pem")
	if err != nil {
		return err
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return err
	}
	jwtToken, err = jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return err
	} else if _, ok := jwtToken.Claims.(*MyCustomClaims); ok {
		return nil
	} else {
		return err
	}

}
