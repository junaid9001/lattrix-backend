package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// migrated
type API struct {
	ID          uuid.UUID `gorm:"primaryKey;not null;type:uuid"`
	UserID      uint      `gorm:"not null"`
	ApiGroupID  uuid.UUID `gorm:"type:uuid;not null"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"size:50;not null"`
	Description *string   `gorm:"size:255"`

	URL    string `gorm:"type:text;not null"`
	Method string `gorm:"size:20;not null;default:GET"`

	AuthType  string  `gorm:"size:20;not null"` //none ,bearer,api-key
	AuthIn    *string `gorm:"size:20"`          //header/query
	AuthKey   *string `gorm:"size:100"`         //authorization,xapikey
	AuthValue *string `gorm:"type:text"`

	Headers  datatypes.JSON `gorm:"type:jsonb"`
	BodyType *string        `gorm:"size:30"` //json,form-data,none
	Body     datatypes.JSON `gorm:"type:jsonb"`

	IntervalSeconds int  `gorm:"not null;default:60"`
	TimeoutMs       int  `gorm:"not null;default:3000"`
	IsActive        bool `gorm:"default:true;index"`

	ExpectedStatusCodes    datatypes.JSON `gorm:"type:jsonb"` //[200,201]
	ExpectedResponseTimeMs *int
	ExpectedBodyContains   *string `gorm:"size:20"`

	LastCheckedAt *time.Time `gorm:"index"`
	LastStatus    string     `gorm:"size:20"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
