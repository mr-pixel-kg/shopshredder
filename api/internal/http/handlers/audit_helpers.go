package handlers

import (
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
	auditcontracts "github.com/mr-pixel-kg/shopshredder/api/internal/auditlog"
	mw "github.com/mr-pixel-kg/shopshredder/api/internal/http/middleware"
	"github.com/mr-pixel-kg/shopshredder/api/internal/services"
)

func newAuditLogInput(
	r *http.Request,
	userID *uuid.UUID,
	action auditcontracts.Action,
	resourceType *auditcontracts.ResourceType,
	resourceID *uuid.UUID,
	details map[string]any,
) services.AuditLogInput {
	return services.AuditLogInput{
		Actor: services.AuditActor{
			UserID:    userID,
			IPAddress: optionalString(extractIP(r)),
			UserAgent: optionalString(strings.TrimSpace(r.UserAgent())),
			ClientID:  mw.ClientIDFromContext(r),
		},
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
	}
}

func newAuditActor(r *http.Request, userID *uuid.UUID) services.AuditActor {
	return newAuditLogInput(r, userID, "", nil, nil, nil).Actor
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func extractIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if i := strings.IndexByte(fwd, ','); i > 0 {
			return strings.TrimSpace(fwd[:i])
		}
		return strings.TrimSpace(fwd)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
