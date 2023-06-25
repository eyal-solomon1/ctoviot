package token

import (
	"testing"
	"time"

	"github.com/eyal-solomon1/ctoviot/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	secretKey := util.RandomString(32)
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	duration := time.Minute

	token, initPayload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, initPayload)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuesAt, time.Now(), time.Second)
	require.WithinDuration(t, payload.ExpiredAt, time.Now().Add(time.Minute), time.Second)

}

func TestExpiredJWTToken(t *testing.T) {
	secretKey := util.RandomString(32)
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	duration := -time.Minute

	token, initPayload, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, initPayload)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, errExpiredToken.Error())
	require.Nil(t, payload)

}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	username := util.RandomOwner()
	duration := time.Minute
	payload, err := NewPayload(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	secretKey := util.RandomString(32)
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, errInvalidToken.Error())
	require.Nil(t, payload)

}
