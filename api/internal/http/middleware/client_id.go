package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const clientContextKey = "client"

const (
	clientIDHeaderName = "X-Client-Id"
	clientIDCookieName = "client_id"
)

func EnsureClientID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(clientContextKey, extractClientID(c))
			return next(c)
		}
	}
}

func ClientID(c echo.Context) *uuid.UUID {
	clientID, ok := c.Get(clientContextKey).(*uuid.UUID)
	if !ok {
		return nil
	}
	return clientID
}

func extractClientID(c echo.Context) *uuid.UUID {
	if id := parseClientID(c.Request().Header.Get(clientIDHeaderName)); id != nil {
		return id
	}

	if cookie, err := c.Cookie(clientIDCookieName); err == nil && cookie != nil {
		return parseClientID(cookie.Value)
	}

	return nil
}

func parseClientID(raw string) *uuid.UUID {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return nil
	}
	return &id
}
