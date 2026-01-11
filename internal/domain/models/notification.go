package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uint      `gorm:"not null;index"`
	Type        string    `gorm:"not null;index"` // workspace_invite
	ReferenceID uuid.UUID `gorm:"type:uuid;not null"`
	ReadAt      *time.Time
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
