package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/manuel/shopware-testenv-platform/api/internal/config"
	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"github.com/manuel/shopware-testenv-platform/api/internal/repositories"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailNotWhitelisted = errors.New("email not whitelisted")
)

type AuthService struct {
	users     *repositories.UserRepository
	passwords *PasswordService
	tokens    *TokenService
	regCfg    config.RegistrationConfig
}

func NewAuthService(
	users *repositories.UserRepository,
	passwords *PasswordService,
	tokens *TokenService,
	regCfg config.RegistrationConfig,
) *AuthService {
	return &AuthService{
		users:     users,
		passwords: passwords,
		tokens:    tokens,
		regCfg:    regCfg,
	}
}

func (s *AuthService) Register(email, password string) (*models.User, error) {
	passwordHash, err := s.passwords.Hash(password)
	if err != nil {
		return nil, err
	}

	if s.regCfg.Mode == config.RegistrationModeWhitelist {
		user, err := s.users.FindPendingByEmail(email)
		if err != nil {
			return nil, ErrEmailNotWhitelisted
		}
		user.PasswordHash = passwordHash
		if err := s.users.Update(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	role := models.RoleUser
	if s.regCfg.AutoAdmin {
		role = models.RoleAdmin
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	}

	if err := s.users.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	user, err := s.users.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if user.IsPending() {
		return "", nil, ErrInvalidCredentials
	}

	if err := s.passwords.Verify(user.PasswordHash, password); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, _, err := s.tokens.Generate(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) Authenticate(tokenValue string) (*models.User, error) {
	claims, err := s.tokens.Parse(tokenValue)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}

	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Logout() error {
	return nil
}
