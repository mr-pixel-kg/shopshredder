package repositories

import (
	"github.com/google/uuid"
	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"gorm.io/gorm"
)

type SandboxEventRepository struct {
	db *gorm.DB
}

func NewSandboxEventRepository(db *gorm.DB) *SandboxEventRepository {
	return &SandboxEventRepository{db: db}
}

func (r *SandboxEventRepository) Create(event *models.SandboxEvent) error {
	return r.db.Create(event).Error
}

func (r *SandboxEventRepository) ListBySandboxID(sandboxID uuid.UUID) ([]models.SandboxEvent, error) {
	var events []models.SandboxEvent
	err := r.db.Where("sandbox_id = ?", sandboxID).Order("created_at asc").Find(&events).Error
	return events, err
}
