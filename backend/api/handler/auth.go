package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/middleware"
	"net/http"
	"strconv"
)

// Auth godoc
// @Summary Endpoint to test the authentication.
// @Description Test the authentication.
// @Tags System Management
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/auth [get]
func AuthCheckHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]interface{}{
		"loggedIn": strconv.FormatBool(middleware.IsUserLoggedIn(c)),
		"username": middleware.GetCurrentUserName(c),
	})

}
