package repositories

import (
	"github.com/google/uuid"
	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) ListPublic() ([]models.Image, error) {
	var images []models.Image
	err := r.db.Where("is_public = ? AND status = ?", true, "ready").Order("created_at desc").Find(&images).Error
	return images, err
}

func (r *ImageRepository) ListAll() ([]models.Image, error) {
	var images []models.Image
	err := r.db.Order("created_at desc").Find(&images).Error
	return images, err
}

func (r *ImageRepository) FindByID(id uuid.UUID) (*models.Image, error) {
	var image models.Image
	if err := r.db.First(&image, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *ImageRepository) Create(image *models.Image) error {
	return r.db.Create(image).Error
}

func (r *ImageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Image{}, "id = ?", id).Error
}

func (r *ImageRepository) UpdateStatus(id uuid.UUID, status string, errorMessage *string) error {
	return r.db.Model(&models.Image{}).Where("id = ?", id).Updates(map[string]any{
		"status":        status,
		"error_message": errorMessage,
	}).Error
}

func (r *ImageRepository) ResetStalePulls(fromStatus, toStatus string, errorMessage *string) error {
	return r.db.Model(&models.Image{}).
		Where("status = ?", fromStatus).
		Updates(map[string]any{"status": toStatus, "error_message": errorMessage}).Error
}
