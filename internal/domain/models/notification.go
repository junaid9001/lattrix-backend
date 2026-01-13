package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes" // Make sure to import this
	"gorm.io/gorm"
)

type Notification struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Type        string    `gorm:"not null;index" json:"type"` // e.g., "invitation"
	Title       string    `gorm:"size:100" json:"title"`
	Message     string    `gorm:"type:text" json:"message"`
	ReferenceID uuid.UUID `gorm:"type:uuid;not null" json:"reference_id"`

	Data datatypes.JSON `gorm:"type:jsonb" json:"data"`

	IsRead    bool           `gorm:"default:false" json:"is_read"`
	ReadAt    *time.Time     `json:"read_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
