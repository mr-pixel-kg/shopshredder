package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      *uuid.UUID `gorm:"type:uuid;index" json:"userId,omitempty"`
	SessionType string     `gorm:"size:32;not null;index" json:"sessionType"`
	TokenID     string     `gorm:"size:255;not null;uniqueIndex" json:"tokenId"`
	ExpiresAt   time.Time  `gorm:"not null;index" json:"expiresAt"`
	BaseModel
}

func (Session) TableName() string {
	return "sessions"
}
