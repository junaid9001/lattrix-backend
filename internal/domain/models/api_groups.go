package models

import (
	"time"

	"github.com/google/uuid"
)

// migrated
type ApiGroup struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;not null;index"`
	WorkspaceID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Name           string    `gorm:"size:50;not null"`
	CreatedByEmail string    `gorm:"size:255"`
	CreatedByID    uint      `gorm:"not null"`
	Description    string    `gorm:"size:100"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
