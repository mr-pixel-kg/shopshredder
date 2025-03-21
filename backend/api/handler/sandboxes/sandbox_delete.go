package sandboxes

import (
	"github.com/labstack/echo/v4"
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/database/models"
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/middleware"
	"log"
	"net/http"
)

type SandboxDeleteResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Docker Image removed successfully"`
}

// SandboxDeleteHandler removes a docker sandbox
// @Summary Remove Sandbox Environment
// @Description Removes a sandbox docker container
// @Tags Sandbox Management
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID"
// @Success 200 {object} SandboxDeleteResponse
// @Failure 400 {object} map[string]string
// @Security BasicAuth
// @Router /api/sandboxes/{id} [delete]
func (h *SandboxHandler) SandboxDeleteHandler(c echo.Context) error {

	ctx := c.Request().Context()
	sandboxId := c.Param("id")

	// Check if sandbox belongs to session or user is logged in
	sessions := h.GuardService.GetSessions(c.RealIP())
	found := middleware.IsUserLoggedIn(c)
	for _, s := range sessions {
		if s.SandboxID == sandboxId {
			found = true
			break
		}
	}
	if !found {
		return echo.NewHTTPError(http.StatusForbidden, "This sandbox does not belong to your session!")
	}

	h.SandboxService.DeleteSandbox(ctx, sandboxId)

	// Write audit log
	h.AuditLogService.LogRequest(c, models.SANDBOX_DELETE, map[string]interface{}{
		"sandbox_id": sandboxId,
	})

	// Remove sandbox session
	err := h.GuardService.UnregisterSession(sandboxId)
	if err != nil {
		log.Printf("Failed to remove sandbox session: %v", err)
	}

	output := SandboxDeleteResponse{
		Message: "Sandbox " + sandboxId + " removed successfully",
		Status:  "success",
	}

	return c.JSON(http.StatusOK, output)
}
