package models

import (
	"time"

	"github.com/google/uuid"
)

type ApiCheckResult struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	APIID     uuid.UUID `gorm:"type:uuid;not null;index"`
	ApiName   string    `json:"api_name"`
	CheckedAt time.Time `gorm:"not null;index" json:"checked_at"`
	Status    string    `gorm:"not null;size:20" json:"status"` //up/down

	StatusCode     *int `json:"status_code"`
	ResponseTimeMs *int `json:"response_time_ms"`

	ErrorMessage *string `json:"error_message"`

	Success bool `gorm:"not null" json:"success"`

	DnsMs        int `json:"dns_ms"`        // Domain Lookup
	TcpMs        int `json:"tcp_ms"`        // Connection Establishment
	TlsMs        int `json:"tls_ms"`        // SSL Handshake
	ProcessingMs int `json:"processing_ms"` // Time To First Byte (Server "Think" Time)
	TransferMs   int `json:"transfer_ms"`   // Content Download Time

	// --- SSL ---
	SslExpiry        *time.Time `json:"ssl_expiry"`
	SslDaysRemaining *int       `json:"ssl_days_remaining"`

	IsStatusCodeMatch   bool `json:"is_status_code_match"`
	IsBodyMatch         bool `json:"is_body_match"`
	IsResponseTimeMatch bool `json:"is_response_time_match"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
