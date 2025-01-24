package images

import (
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/services"
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/services/images"
)

type ImageHandler struct {
	DockerService *services.DockerService
	ImageService  *images.ImageService
}

func NewImageHandler(dockerService *services.DockerService, imageService *images.ImageService) *ImageHandler {
	return &ImageHandler{
		DockerService: dockerService,
		ImageService:  imageService,
	}
}
