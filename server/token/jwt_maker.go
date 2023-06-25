package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 32

// JWTMaker is a JWT maker
type JWTMaker struct {
	secretKey string
}

func (maker JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// meaning the token header 'alg' doesn't match our signing algorithem
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errInvalidToken) {
			return nil, errInvalidToken
		}
		return nil, errExpiredToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errInvalidToken
	}

	return payload, nil

}

// CreateToken creates a new token based on the 'username' input
func (maker JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))

	return token, payload, err
}

// NewJWTMaker creates a new JWTmaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, errors.New(fmt.Sprintf("invalid key size:must be at least %d", minSecretKeySize))

	}
	return JWTMaker{secretKey}, nil
}
