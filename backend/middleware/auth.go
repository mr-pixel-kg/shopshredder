package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/config"
	"log"
	"log/slog"
	"strconv"
)

func AuthRequiredMiddleware(authConfig config.AuthConfig) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		c.Response().Header().Set("X-Auth", "required")

		// Check password
		if username != authConfig.Username || password != authConfig.Password {
			log.Printf("Wrong login credentials for user %s\n", username)
			return false, nil // Wrong credentials
		}

		// Store user info
		ctx := c.(*AuthContext)
		ctx.IsAuthenticated = true
		ctx.Username = username
		c = ctx
		return true, nil
	})
}

type AuthContext struct {
	echo.Context
	IsAuthenticated bool
	Username        string
}

func AuthMiddleware(authConfig config.AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			c := &AuthContext{ctx, false, ""}

			username, password, ok := c.Request().BasicAuth()
			if ok == true && username == authConfig.Username && password == authConfig.Password {
				c.IsAuthenticated = true
				c.Username = username
			} else {
				if username != "" {
					slog.Warn("Request with invalid login credentials")
				}
			}

			ctx.Response().Header().Set("X-Authenticated", strconv.FormatBool(c.IsAuthenticated))
			ctx.Response().Header().Set("X-Username", c.Username)

			return next(c)
		}
	}
}

func IsUserLoggedIn(context echo.Context) bool {
	c := context.(*AuthContext)
	return c.IsAuthenticated
}

func GetCurrentUserName(context echo.Context) *string {
	c := context.(*AuthContext)
	if c.IsAuthenticated {
		return &c.Username
	}
	return nil
}
