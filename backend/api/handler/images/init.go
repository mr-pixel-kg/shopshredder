package images

import (
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/services"
)

type ImageHandler struct {
	DockerService   *services.DockerService
	ImageService    *services.ImageService
	AuditLogService *services.AuditLogService
}

func NewImageHandler(dockerService *services.DockerService, imageService *services.ImageService, auditLogService *services.AuditLogService) *ImageHandler {
	return &ImageHandler{
		DockerService:   dockerService,
		ImageService:    imageService,
		AuditLogService: auditLogService,
	}
}
