package auth

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/vardius/go-api-boilerplate/pkg/errors"
)

type Authenticator interface {
	Verify(token string, claims jwt.Claims) error
	Sign(token *jwt.Token) (string, error)
}

func NewSecretAuthenticator(secret []byte) Authenticator {
	return secretAuthenticator{
		secretKey: secret,
	}
}

type secretAuthenticator struct {
	secretKey []byte
}

func (a secretAuthenticator) Sign(token *jwt.Token) (string, error) {
	return token.SignedString(a.secretKey)
}
func (a secretAuthenticator) Verify(token string, claims jwt.Claims) (err error) {
	accessToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Failed to decode token, invalid signing method")
		}
		return a.secretKey, nil
	})
	if err != nil {
		return errors.New("Failed to parse token with claims")
	}

	if !accessToken.Valid {
		return errors.New("Token is not valid, could not parse claims")
	}

	return nil
}
