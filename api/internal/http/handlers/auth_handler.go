package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
		return responses.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
	}

	user, err := h.auth.Register(input.Email, input.Password)
	if err != nil {
		return responses.Error(c, http.StatusBadRequest, "REGISTER_FAILED", err.Error())
	}

	_ = h.audit.Log(&user.ID, "auth.registered", c.RealIP(), map[string]any{"email": user.Email})
	return c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var input dto.LoginRequest
	if err := c.Bind(&input); err != nil {
		return responses.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
	}

	token, user, err := h.auth.Login(input.Email, input.Password)
	if err != nil {
		return responses.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email or password is invalid")
	}

	_ = h.audit.Log(&user.ID, "auth.logged_in", c.RealIP(), map[string]any{})
	return c.JSON(http.StatusOK, map[string]any{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	auth := mw.MustAuth(c)
	if err := h.auth.Logout(auth.TokenID); err != nil {
		return responses.Error(c, http.StatusInternalServerError, "LOGOUT_FAILED", "Could not log out")
	}

	_ = h.audit.Log(&auth.UserID, "auth.logged_out", c.RealIP(), map[string]any{})
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) Me(c echo.Context) error {
	return c.JSON(http.StatusOK, c.Get("user"))
}
