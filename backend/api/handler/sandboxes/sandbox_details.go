package sandboxes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// SandboxDetailsHandler returns information about a sandbox
// @Summary Details about a Sandbox
// @Description Get details about a sandbox environment
// @Tags Sandbox Management
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID" example(67777b4e-946f-4462-b689-3c608d2d7938)
// @Success 200 {object} sandbox.SandboxInfo
// @Failure 400 {object} map[string]string
// @Router /api/sandboxes/{id} [get]
func (h *SandboxHandler) SandboxDetailsHandler(c echo.Context) error {

	ctx := c.Request().Context()
	sandboxId := c.Param("id")

	sandbox, err := h.SandboxService.GetSandbox(ctx, sandboxId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, sandbox)
}
