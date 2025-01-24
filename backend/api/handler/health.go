package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// HealthCheck godoc
// @Summary Show the status of server.
// @Description Get the status of server.
// @Tags System Management
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/health [get]
func HealthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "OK",
		"message": "API is running smoothly",
	})
}
