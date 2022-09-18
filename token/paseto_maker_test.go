package token

import (
	"testing"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/khorsl/simple_bank/db/util"
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	symmetricKey := util.RandomString(chacha20poly1305.KeySize)

	pasetoMaker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)
	require.NotEmpty(t, pasetoMaker)

	username := util.RandomUsername()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = pasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	symmetricKey := util.RandomString(chacha20poly1305.KeySize)

	pasetoMaker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	username := util.RandomUsername()
	duration := -time.Minute

	token, payload, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = pasetoMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Empty(t, payload)
}

func TestPasetoMakerSymmetricKeySize(t *testing.T) {
	symmetricKey := util.RandomString(chacha20poly1305.KeySize - 1)

	pasetoMaker, err := NewPasetoMaker(symmetricKey)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidPasetoKeySize.Error())
	require.Empty(t, pasetoMaker)
}

func TestInvalidPasetoToken(t *testing.T) {
	symmetricKey := util.RandomString(chacha20poly1305.KeySize)

	pasetoMaker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	// Encrypt an invalid token with another symmetric key
	symmetricKey2 := []byte(util.RandomString(chacha20poly1305.KeySize))
	invalidToken, err := paseto.NewV2().Encrypt(symmetricKey2, &Payload{}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, invalidToken)

	payload, err := pasetoMaker.VerifyToken(invalidToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Empty(t, payload)
}
