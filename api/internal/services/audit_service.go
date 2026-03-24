package services

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"github.com/manuel/shopware-testenv-platform/api/internal/repositories"
	"gorm.io/datatypes"
)

type AuditService struct {
	repo *repositories.AuditLogRepository
}

func NewAuditService(repo *repositories.AuditLogRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Log(userID *uuid.UUID, action, ip string, details map[string]any) error {
	payload, err := json.Marshal(details)
	if err != nil {
		return err
	}

	return s.repo.Create(&models.AuditLog{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		IPAddress: ip,
		Details:   datatypes.JSON(payload),
		CreatedAt: time.Now().UTC(),
	})
}

func (s *AuditService) List(limit int) ([]models.AuditLog, error) {
	return s.repo.List(limit)
}
