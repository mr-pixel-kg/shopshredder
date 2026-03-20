package models

import "github.com/google/uuid"

type Image struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name            string     `gorm:"size:255;not null" json:"name"`
	Tag             string     `gorm:"size:255;not null" json:"tag"`
	Title           *string    `gorm:"size:255" json:"title,omitempty"`
	Description     *string    `gorm:"type:text" json:"description,omitempty"`
	ThumbnailURL    *string    `gorm:"-" json:"thumbnailUrl,omitempty"` /* This field is not stored in the database, it's computed based on the image ID */
	IsPublic        bool       `gorm:"not null;default:false" json:"isPublic"`
	CreatedByUserID *uuid.UUID `gorm:"type:uuid" json:"createdByUserId,omitempty"`
	BaseModel
}

func (Image) TableName() string {
	return "images"
}

func (i Image) FullName() string {
	return i.Name + ":" + i.Tag
}
