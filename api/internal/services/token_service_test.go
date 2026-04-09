package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mr-pixel-kg/shopshredder/api/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenServiceGenerateAndParse(t *testing.T) {
	t.Parallel()

	service := NewTokenService(config.AuthConfig{
		JWTSecret:     "unit-test-secret",
		JWTTTLMinutes: 30,
	})
	userID := uuid.New()

	token, expiresAt, err := service.Generate(userID)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	assert.WithinDuration(t, time.Now().UTC().Add(30*time.Minute), expiresAt, 5*time.Second)

	claims, err := service.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, userID.String(), claims.Subject)
}

func TestTokenServiceParseRejectsTamperedToken(t *testing.T) {
	t.Parallel()

	service := NewTokenService(config.AuthConfig{
		JWTSecret:     "unit-test-secret",
		JWTTTLMinutes: 30,
	})

	token, _, err := service.Generate(uuid.New())
	require.NoError(t, err)

	tampered := token + "tampered"
	claims, err := service.Parse(tampered)
	require.Error(t, err)
	assert.Nil(t, claims)
}

func TestTokenServiceParseRejectsExpiredToken(t *testing.T) {
	t.Parallel()

	service := NewTokenService(config.AuthConfig{JWTSecret: "unit-test-secret"})

	expired := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   uuid.NewString(),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC().Add(-2 * time.Minute)),
	})

	token, err := expired.SignedString([]byte("unit-test-secret"))
	require.NoError(t, err)

	claims, err := service.Parse(token)
	require.Error(t, err)
	assert.Nil(t, claims)
}
