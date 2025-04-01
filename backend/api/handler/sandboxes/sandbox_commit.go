package sandboxes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// SandboxCommitHandler creates a new sandbox image from the current sandbox container
// @Summary Creates a new sandbox image from the sandbox container
// @Description Creates a new sandbox image from the sandbox container
// @Tags Sandbox Management
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID" example(67777b4e-946f-4462-b689-3c608d2d7938)
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security BasicAuth
// @Router /api/sandboxes/{id}/commit [post]
func (h *SandboxHandler) SandboxCommitHandler(c echo.Context) error {

	ctx := c.Request().Context()
	sandboxId := c.Param("id")

	h.SandboxService.CommitSandbox(ctx, sandboxId, "sandbox-commit:latest")
	/*if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}*/

	return c.JSON(http.StatusOK, map[string]string{
		"message":   "Sandbox committed successfully",
		"sandboxId": sandboxId,
	})
}
