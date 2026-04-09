package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mr-pixel-kg/shopshredder/api/internal/config"
)

type TokenService struct {
	cfg config.AuthConfig
}

func NewTokenService(cfg config.AuthConfig) *TokenService {
	return &TokenService{cfg: cfg}
}

func (s *TokenService) Generate(userID uuid.UUID) (string, time.Time, error) {
	expiresAt := time.Now().UTC().Add(time.Duration(s.cfg.JWTTTLMinutes) * time.Minute)

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	return signed, expiresAt, err
}

func (s *TokenService) Parse(tokenValue string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
