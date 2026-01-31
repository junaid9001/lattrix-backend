package repository

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"gorm.io/gorm"
)

type ApiGroupRepository interface {
	Create(*models.ApiGroup) error
	Delete(ID, workspaceID uuid.UUID) error
	FindByID(ID, workspaceID uuid.UUID) (*models.ApiGroup, error)
	Update(id, workspaceID uuid.UUID, updates map[string]any) (*models.ApiGroup, error)
	WithDB(db *gorm.DB) ApiGroupRepository
	ListGroups(workspaceID uuid.UUID) (*[]models.ApiGroup, error)
}
