package services

import (
	"github.com/google/uuid"
	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"github.com/manuel/shopware-testenv-platform/api/internal/repositories"
)

type ImageService struct {
	repo *repositories.ImageRepository
}

func NewImageService(repo *repositories.ImageRepository) *ImageService {
	return &ImageService{repo: repo}
}

func (s *ImageService) ListPublic() ([]models.Image, error) {
	return s.repo.ListPublic()
}

func (s *ImageService) ListAll() ([]models.Image, error) {
	return s.repo.ListAll()
}

func (s *ImageService) Create(image *models.Image) error {
	image.ID = uuid.New()
	return s.repo.Create(image)
}

func (s *ImageService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *ImageService) FindByID(id uuid.UUID) (*models.Image, error) {
	return s.repo.FindByID(id)
}
