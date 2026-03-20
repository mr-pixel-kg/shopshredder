package types

import "github.com/google/uuid"

type AuthContext struct {
	UserID      uuid.UUID
	TokenID     string
	SessionType string
}
