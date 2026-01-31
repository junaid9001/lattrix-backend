package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// migrated
type API struct {
	ID          uuid.UUID `gorm:"primaryKey;not null;type:uuid" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	ApiGroupID  uuid.UUID `gorm:"type:uuid;not null" json:"api_group_id"`
	WorkspaceID uuid.UUID `gorm:"type:uuid;not null" json:"workspace_id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Description *string   `gorm:"size:255" json:"description"`

	URL    string `gorm:"type:text;not null" json:"url"`
	Method string `gorm:"size:20;not null;default:GET" json:"method"`

	AuthType  string  `gorm:"size:20;not null" json:"auth_type"` //none ,bearer,api-key
	AuthIn    *string `gorm:"size:20" json:"auth_in"`            //header/query
	AuthKey   *string `gorm:"size:100" json:"auth_key"`          //authorization,xapikey
	AuthValue *string `gorm:"type:text" json:"auth_value"`

	Headers  datatypes.JSON `gorm:"type:jsonb" json:"headers"`
	BodyType *string        `gorm:"size:30" json:"body_type"` //json,form-data,none
	Body     datatypes.JSON `gorm:"type:jsonb" json:"body"`

	IntervalSeconds int  `gorm:"not null;default:60" json:"interval_seconds"`
	TimeoutMs       int  `gorm:"not null;default:3000" json:"timeout_ms"`
	IsActive        bool `gorm:"default:true;index" json:"is_active"`

	ExpectedStatusCodes    datatypes.JSON `gorm:"type:jsonb" json:"expected_status_codes"` //[200,201]
	ExpectedResponseTimeMs *int           `json:"expected_response_time_ms"`
	ExpectedBodyContains   *string        `gorm:"size:200" json:"expected_body_contains"`

	LastCheckedAt *time.Time `gorm:"index" json:"last_checked_at"`
	NextCheckAt   time.Time  `gorm:"not null;index" json:"next_check_at"`
	LastStatus    string     `gorm:"size:20" json:"last_status"`

	LastResponseTimeMs int     `gorm:"default:0" json:"last_response_time_ms"`
	LastErrorMessage   *string `gorm:"size:255" json:"last_error_message"`

	NotifyAfterFailures int `gorm:"not null;default:1" json:"notify_after_failures"`
	DownCount           int `gorm:"not null;default:0" json:"down_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
