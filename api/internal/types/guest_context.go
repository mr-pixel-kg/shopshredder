package types

import "github.com/google/uuid"

type GuestContext struct {
	SessionID uuid.UUID
	TokenID   string
}
