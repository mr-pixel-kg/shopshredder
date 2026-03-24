package repositories

import (
	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(entry *models.AuditLog) error {
	return r.db.Create(entry).Error
}

func (r *AuditLogRepository) List(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email")
		}).
		Order("created_at desc").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}
