package handlers

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	auditcontracts "github.com/manuel/shopware-testenv-platform/api/internal/auditlog"
	"github.com/manuel/shopware-testenv-platform/api/internal/http/dto"
	mw "github.com/manuel/shopware-testenv-platform/api/internal/http/middleware"
	"github.com/manuel/shopware-testenv-platform/api/internal/http/responses"
	"github.com/manuel/shopware-testenv-platform/api/internal/logging"
	"github.com/manuel/shopware-testenv-platform/api/internal/services"
)

type WhitelistHandler struct {
	users *services.UserService
	audit *services.AuditService
}

func NewWhitelistHandler(users *services.UserService, audit *services.AuditService) *WhitelistHandler {
	return &WhitelistHandler{users: users, audit: audit}
}

// List godoc
// @Summary      List whitelisted emails
// @Description  Return all pending (whitelisted but not yet registered) users
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} dto.UserResponse
// @Failure      403 {object} dto.ErrorResponse
// @Router       /api/whitelist [get]
func (h *WhitelistHandler) List(c echo.Context) error {
	users, err := h.users.ListPending()
	if err != nil {
		return responses.FromError(c, err)
	}
	out := make([]dto.UserResponse, len(users))
	for i, u := range users {
		out[i] = dto.UserResponse{
			ID: u.ID, Email: u.Email, Role: u.Role,
			IsPending: u.IsPending(), CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
		}
	}
	return c.JSON(200, out)
}

// Add godoc
// @Summary      Add email to whitelist
// @Description  Create a pending user row so the email can register in whitelist mode
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.AddWhitelistRequest true "Email to whitelist"
// @Success      201 {object} dto.UserResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      409 {object} dto.ErrorResponse
// @Router       /api/whitelist [post]
func (h *WhitelistHandler) Add(c echo.Context) error {
	var input dto.AddWhitelistRequest
	if err := bindAndValidate(c, &input); err != nil {
		return responses.FromError(c, err)
	}

	auth := mw.MustAuth(c)
	user, err := h.users.AddWhitelist(input.Email, input.Role)
	if err != nil {
		return responses.FromError(c, err)
	}

	slog.Info("email whitelisted", logging.RequestFields(c,
		"component", "admin",
		"whitelisted_email", logging.MaskEmail(input.Email),
	)...)
	resourceType := auditcontracts.ResourceTypeUser
	_ = h.audit.Log(newAuditLogInput(c, &auth.UserID, auditcontracts.ActionUserWhitelisted, &resourceType, &user.ID, map[string]any{
		"email": user.Email,
		"role":  user.Role,
	}))
	return c.JSON(201, dto.UserResponse{
		ID: user.ID, Email: user.Email, Role: user.Role,
		IsPending: user.IsPending(), CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
	})
}

// Remove godoc
// @Summary      Remove email from whitelist
// @Description  Delete a pending user row (only works for users that have not yet registered)
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Pending user ID" format(uuid)
// @Success      204
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Router       /api/whitelist/{id} [delete]
func (h *WhitelistHandler) Remove(c echo.Context) error {
	id, err := parseUUIDParam(c, "id", "INVALID_ID", "Invalid UUID")
	if err != nil {
		return responses.FromError(c, err)
	}

	user, getErr := h.users.Get(id)
	if getErr != nil {
		return responses.FromError(c, getErr)
	}

	if err := h.users.RemoveWhitelist(id); err != nil {
		return responses.FromError(c, err)
	}

	auth := mw.MustAuth(c)
	slog.Info("email removed from whitelist", logging.RequestFields(c,
		"component", "admin",
		"removed_email", logging.MaskEmail(user.Email),
	)...)
	resourceType := auditcontracts.ResourceTypeUser
	_ = h.audit.Log(newAuditLogInput(c, &auth.UserID, auditcontracts.ActionUserWhitelistRemoved, &resourceType, &user.ID, map[string]any{
		"email": user.Email,
	}))
	return c.NoContent(204)
}
