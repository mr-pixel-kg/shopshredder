package middleware

import (
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/manuel/shopware-testenv-platform/api/internal/apperror"
	"github.com/manuel/shopware-testenv-platform/api/internal/http/responses"
	"github.com/manuel/shopware-testenv-platform/api/internal/logging"
	"github.com/manuel/shopware-testenv-platform/api/internal/models"
	"github.com/manuel/shopware-testenv-platform/api/internal/services"
	"github.com/manuel/shopware-testenv-platform/api/internal/types"
)

const authContextKey = "auth"

func Auth(authService *services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := strings.TrimSpace(c.Request().Header.Get(echo.HeaderAuthorization))
			if authHeader == "" {
				slog.Warn("missing authorization header", logging.RequestFields(c, "component", "auth")...)
				return responses.FromAppError(c, apperror.Unauthorized("Missing bearer token"))
			}

			token, ok := ParseAuthorizationHeader(authHeader)
			if !ok {
				slog.Warn("invalid authorization header format", logging.RequestFields(c, "component", "auth")...)
				return responses.FromAppError(c, apperror.Unauthorized("Invalid authorization header"))
			}

			user, err := authService.Authenticate(token)
			if err != nil {
				slog.Warn("token authentication failed", append(logging.RequestFields(c, "component", "auth"), "error", err.Error())...)
				return responses.FromAppError(c, apperror.Unauthorized("Invalid or expired token"))
			}

			c.Set(authContextKey, types.AuthContext{UserID: user.ID})
			c.Set("user", user)
			slog.Debug("request authenticated", logging.RequestFields(c, "component", "auth", "user_id", user.ID.String())...)
			return next(c)
		}
	}
}

func MustAuth(c echo.Context) types.AuthContext {
	return c.Get(authContextKey).(types.AuthContext)
}

func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*models.User)
			if !ok || !user.IsAdmin() {
				return responses.FromAppError(c, apperror.Forbidden("Admin access required"))
			}
			return next(c)
		}
	}
}

func ParseAuthorizationHeader(authHeader string) (string, bool) {
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}
