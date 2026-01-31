package models

import (
	"time"

	"github.com/google/uuid"
)

// migrated
type ApiGroup struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;not null;index" json:"id"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;not null;index" json:"workspace_id"`

	Name           string    `gorm:"size:50;not null" json:"name"`
	CreatedByEmail string    `gorm:"size:255" json:"created_by_email"`
	CreatedByID    uint      `gorm:"not null" json:"created_by_id"`
	Description    string    `gorm:"size:100" json:"description"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
