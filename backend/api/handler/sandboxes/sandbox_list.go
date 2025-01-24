package sandboxes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// List all sandbox environments
// @Summary List sandbox environments
// @Description List sandbox environments
// @Tags Sandbox Management
// @Accept json
// @Produce json
// @Success 200 {object} sandbox.SandboxInfo
// @Failure 400 {object} map[string]string
// @Router /api/sandboxes [get]
func (h *SandboxHandler) SandboxListHandler(c echo.Context) error {
	ctx := c.Request().Context()

	sandboxes, err := h.SandboxService.ListSandboxes(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, sandboxes)
}
