package sandboxes

import (
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/services"
)

type SandboxHandler struct {
	SandboxService  *services.SandboxService
	AuditLogService *services.AuditLogService
}

func NewSandboxHandler(sandboxService *services.SandboxService, auditLogService *services.AuditLogService) *SandboxHandler {
	return &SandboxHandler{
		SandboxService:  sandboxService,
		AuditLogService: auditLogService,
	}
}
