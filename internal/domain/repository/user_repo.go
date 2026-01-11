package repository

//what operation the business needs from user storage

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(*models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(ID int) (*models.User, error)
	UpdateProfile(userID int, updates map[string]interface{}) (*models.User, error)
	CreateWorkSpace(userID uint) (uuid.UUID, error)
	WithDB(db *gorm.DB) UserRepository
}
