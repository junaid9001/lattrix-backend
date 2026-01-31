package models

//whats user
import (
	"gorm.io/gorm"
)

type PlanType string

const (
	PlanFree   PlanType = "FREE"
	PlanPro    PlanType = "PRO"
	PlanAgency PlanType = "AGENCY"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:100"`
	Email    string `gorm:"size:100;uniqueIndex"`
	Password string `gorm:"size:100"`

	Plan               PlanType `gorm:"default:'FREE'" json:"plan"`
	StripeCustomerID   *string  `json:"stripe_customer_id"`
	SubscriptionStatus string   `gorm:"default:inactive" json:"subscription_status"`

	IsSuperAdmin bool `gorm:"default:false"`

	IsActive      bool `gorm:"default:true"`
	EmailVerified bool `gorm:"default:false"`
}
