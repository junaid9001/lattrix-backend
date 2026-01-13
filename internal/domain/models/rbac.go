package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//all migrated

// permissions are seeded
type Permission struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	Code        string    `gorm:"uniqueIndex;not null" json:"code"` //permission name
	Description string    `json:"description"`
}

type Role struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	WorkspaceID uuid.UUID      `gorm:"type:uuid;not null;index" json:"workspace_id"`
	Name        string         `gorm:"size:50" json:"name"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type RolePermission struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null;index"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null"`
}

type UserRole struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uint      `gorm:"not null;index"`
	RoleID      uuid.UUID `gorm:"type:uuid;index;not null"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;index;not null"`
}
