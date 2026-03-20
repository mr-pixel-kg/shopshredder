package repositories

import (
	"time"

	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) FindByTokenID(tokenID string) (*models.Session, error) {
	var session models.Session
	if err := r.db.First(&session, "token_id = ?", tokenID).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteByTokenID(tokenID string) error {
	return r.db.Where("token_id = ?", tokenID).Delete(&models.Session{}).Error
}

func (r *SessionRepository) ExistsActiveToken(tokenID string, now time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.Session{}).
		Where("token_id = ?", tokenID).
		Where("expires_at > ?", now).
		Count(&count).Error
	return count > 0, err
}
