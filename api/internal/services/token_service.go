package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/manuel/shopware-testenv-platform/api/internal/config"
)

type TokenService struct {
	cfg config.AuthConfig
}

type Claims struct {
	UserID      string `json:"userId,omitempty"`
	SessionType string `json:"sessionType"`
	TokenID     string `json:"tokenId"`
	jwt.RegisteredClaims
}

func NewTokenService(cfg config.AuthConfig) *TokenService {
	return &TokenService{cfg: cfg}
}

func (s *TokenService) Generate(userID uuid.UUID) (string, string, time.Time, error) {
	return s.generate(userID.String(), "user", s.cfg.JWTTTLMinutes)
}

func (s *TokenService) GenerateGuest(sessionID uuid.UUID) (string, string, time.Time, error) {
	return s.generate(sessionID.String(), "guest", s.cfg.GuestJWTTTLMinutes)
}

func (s *TokenService) generate(subject, sessionType string, ttlMinutes int) (string, string, time.Time, error) {
	tokenID := uuid.NewString()
	expiresAt := time.Now().UTC().Add(time.Duration(ttlMinutes) * time.Minute)

	claims := Claims{
		SessionType: sessionType,
		TokenID:     tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	if sessionType == "user" {
		claims.UserID = subject
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	return signed, tokenID, expiresAt, err
}

func (s *TokenService) Parse(tokenValue string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
