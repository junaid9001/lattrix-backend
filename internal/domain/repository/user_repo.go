package repository

//what operation the business needs from user storage

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/http/dto"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(*models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(ID int) (*models.User, error)
	UpdateProfile(userID uint, updates map[string]interface{}) (*models.User, error)
	CreateWorkSpace(userID uint, name string) (uuid.UUID, error)
	WithDB(db *gorm.DB) UserRepository
	WorkspaceUsers(workspaceID uuid.UUID) ([]dto.WorkspaceUsers, error)
	GetUserWorkspaces(userID uint) ([]dto.UserWorkspaceResponse, error)
}
