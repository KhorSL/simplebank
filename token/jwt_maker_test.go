package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/khorsl/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestJwtMaker(t *testing.T) {
	secretKey := util.RandomString(32)

	jwtMaker, err := NewJwtMaker(secretKey)
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	username := util.RandomUsername()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = jwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJwtToken(t *testing.T) {
	secretKey := util.RandomString(32)

	jwtMaker, err := NewJwtMaker(secretKey)
	require.NoError(t, err)

	username := util.RandomUsername()
	duration := -time.Minute

	token, payload, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Empty(t, payload)
}

func TestJwtMakerSecretKeySize(t *testing.T) {
	secretKey := util.RandomString(16)

	jwtMaker, err := NewJwtMaker(secretKey)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidJwtKeySize.Error())
	require.Empty(t, jwtMaker)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	username := util.RandomUsername()
	payload, err := NewPayload(username, time.Minute)
	require.NoError(t, err)

	// Create token with different signing method
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	secretKey := util.RandomString(32)
	jwtMaker, err := NewJwtMaker(secretKey)
	require.NoError(t, err)

	payload, err = jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
