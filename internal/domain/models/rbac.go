package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//all migrated

// permissions are seeded
type Permission struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;"`
	Code        string    `gorm:"uniqueIndex;not null"` //permission name
	Description string
}

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"size:50"` //flexible
	gorm.DeletedAt
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
