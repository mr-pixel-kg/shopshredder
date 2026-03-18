package responses

import "github.com/labstack/echo/v4"

func Error(c echo.Context, status int, code, message string) error {
	return c.JSON(status, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}
