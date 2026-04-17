package models

import "time"

type VerificationOTP struct {
	UserID    uint      `gorm:"primaryKey"`
	Code      string    `gorm:"size:6;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
