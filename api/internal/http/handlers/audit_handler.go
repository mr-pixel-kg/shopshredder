package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/manuel/shopware-testenv-platform/api/internal/http/responses"
	"github.com/manuel/shopware-testenv-platform/api/internal/services"
)

type AuditHandler struct {
	audit *services.AuditService
}

func NewAuditHandler(audit *services.AuditService) *AuditHandler {
	return &AuditHandler{audit: audit}
}

func (h *AuditHandler) List(c echo.Context) error {
	limit := 50
	if value := c.QueryParam("limit"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}

	logs, err := h.audit.List(limit)
	if err != nil {
		return responses.Error(c, http.StatusInternalServerError, "AUDIT_LOG_LIST_FAILED", "Could not load audit logs")
	}
	return c.JSON(http.StatusOK, logs)
}
