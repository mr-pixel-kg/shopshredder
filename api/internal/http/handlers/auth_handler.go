package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/manuel/shopware-testenv-platform/api/internal/apperror"
	"github.com/manuel/shopware-testenv-platform/api/internal/http/dto"
	mw "github.com/manuel/shopware-testenv-platform/api/internal/http/middleware"
	"github.com/manuel/shopware-testenv-platform/api/internal/http/responses"
	"github.com/manuel/shopware-testenv-platform/api/internal/services"
)

type AuthHandler struct {
	auth  *services.AuthService
	audit *services.AuditService
}

func NewAuthHandler(auth *services.AuthService, audit *services.AuditService) *AuthHandler {
	return &AuthHandler{auth: auth, audit: audit}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var input dto.RegisterRequest
	if err := c.Bind(&input); err != nil {
		return responses.FromAppError(c, apperror.BadRequest("VALIDATION_ERROR", "Invalid request body"))
	}

	user, err := h.auth.Register(input.Email, input.Password)
	if err != nil {
		return responses.FromAppError(c, apperror.BadRequest("REGISTER_FAILED", "Could not register user").WithCause(err))
	}

	_ = h.audit.Log(&user.ID, "auth.registered", c.RealIP(), map[string]any{"email": user.Email})
	return c.JSON(201, user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var input dto.LoginRequest
	if err := c.Bind(&input); err != nil {
		return responses.FromAppError(c, apperror.BadRequest("VALIDATION_ERROR", "Invalid request body"))
	}

	token, user, err := h.auth.Login(input.Email, input.Password)
	if err != nil {
		return responses.FromAppError(c, apperror.Unauthorized("Email or password is invalid").WithCause(err))
	}

	_ = h.audit.Log(&user.ID, "auth.logged_in", c.RealIP(), map[string]any{})
	return c.JSON(200, map[string]any{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	auth := mw.MustAuth(c)
	if err := h.auth.Logout(auth.TokenID); err != nil {
		return responses.FromAppError(c, apperror.Internal("LOGOUT_FAILED", "Could not log out").WithCause(err))
	}

	_ = h.audit.Log(&auth.UserID, "auth.logged_out", c.RealIP(), map[string]any{})
	return c.NoContent(204)
}

func (h *AuthHandler) Me(c echo.Context) error {
	return c.JSON(200, c.Get("user"))
}
