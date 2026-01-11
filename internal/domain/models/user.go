package models

//whats user
import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:100"`
	Email    string `gorm:"size:100;uniqueIndex"`
	Password string `gorm:"size:100"`

	Role          string    `gorm:"size:20;not null"`          //superadmin/admin/user
	WorkspaceID   uuid.UUID `gorm:"type:uuid;index;not null;"` //uuid.New() creates new uuid
	IsActive      bool      `gorm:"default:true"`
	EmailVerified bool      `gorm:"default:false"`
}
