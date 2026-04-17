package seeds

import (
	"errors"

	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedSysAdmin(db *gorm.DB, sysAdminEmail, sysAdminPassword string) error {

	result := db.Where("email = ?", sysAdminEmail).First(&models.User{})
	if result.Error == gorm.ErrRecordNotFound {
		pass, _ := bcrypt.GenerateFromPassword([]byte(sysAdminPassword), bcrypt.DefaultCost)
		user := models.User{
			Username:           "sys_admin",
			Email:              sysAdminEmail,
			Password:           string(pass),
			Plan:               "AGENCY",
			SubscriptionStatus: "active",
			IsSuperAdmin:       true,
			IsActive:           true,
			EmailVerified:      false,
		}
		err := db.Create(&user).Error
		if err != nil {
			return errors.New("sys admin config failed")
		}
	}

	return nil
}
