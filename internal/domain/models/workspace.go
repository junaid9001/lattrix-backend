package models

import (
	"time"

	"github.com/google/uuid"
)

// migrated
type Workspace struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerID   uint
	CreatedAt time.Time
}
