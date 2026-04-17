package infra

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/infra/seeds"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	_ = godotenv.Load()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&models.User{}, &models.ApiGroup{}, &models.API{}, &models.Permission{},
		&models.Workspace{}, &models.Role{}, &models.RolePermission{}, &models.UserRole{}, &models.WorkspaceInvitation{},
		&models.Notification{}, &models.ApiCheckResult{}, &models.WorkspaceNotification{}, &models.VerificationOTP{},
	)

	sysAdminEmail := os.Getenv("SYS_ADMIN_EMAIL")
	sysAdminPassword := os.Getenv("SYS_ADMIN_PASSWORD")
	if sysAdminEmail == "" || sysAdminPassword == "" {
		log.Fatal("sys admin configs are missing")
	}

	status := seeds.SeedSysAdmin(db, sysAdminEmail, sysAdminPassword)
	if status != nil {
		log.Fatal("sys admin seeding failed")
	}

	log.Print("db connection success")
	seeds.SeedPermissions(db)
	return db
}
