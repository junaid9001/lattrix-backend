package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkspaceInvitation struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;not null;index"`
	Email       string    `gorm:"not null;index"`
	RoleID      uuid.UUID `gorm:"type:uuid;not null"`
	InvitedBy   uint      `gorm:"not null"`
	Token       string    `gorm:"not null;uniqueIndex"`
	Status      string    `gorm:"not null;index"` // pending, accepted, rejected
	ExpiresAt   time.Time `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
