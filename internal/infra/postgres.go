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
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&models.User{}, &models.ApiGroup{}, &models.API{}, &models.Permission{},
		&models.Workspace{}, &models.Role{}, &models.RolePermission{}, &models.UserRole{}, &models.WorkspaceInvitation{},
		&models.Notification{},
	)
	log.Print("db connection success")
	seeds.SeedPermissions(db)
	return db
}
