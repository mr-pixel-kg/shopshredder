package sandboxes

import (
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/services/sandbox"
)

type SandboxHandler struct {
	SandboxService *sandbox.SandboxService
}

func NewSandboxHandler(sandboxService *sandbox.SandboxService) *SandboxHandler {
	return &SandboxHandler{
		SandboxService: sandboxService,
	}
}
