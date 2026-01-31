package models

import (
	"time"

	"github.com/google/uuid"
)

// migrated
type Workspace struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"size:100;not null;default:'My Workspace'"`
	OwnerID   uint
	CreatedAt time.Time
}

type WorkspaceNotification struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;not null;index" json:"workspace_id"`
	Title       string    `gorm:"size:255" json:"title"`
	Message     string    `gorm:"type:text" json:"message"`
	CreatedAt   time.Time `json:"time"`
}
