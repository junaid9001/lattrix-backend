package repository

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"gorm.io/gorm"
)

type InvitationRepository interface {
	WithTx(tx *gorm.DB) InvitationRepository

	Create(invite *models.WorkspaceInvitation) error
	FindByID(id uuid.UUID) (*models.WorkspaceInvitation, error)
	FindByToken(token string) (*models.WorkspaceInvitation, error)
	UpdateStatus(id uuid.UUID, status string) error
}
