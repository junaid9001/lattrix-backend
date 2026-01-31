package repository

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
)

type WorkspaceNotificationRepository interface {
	Create(notification *models.WorkspaceNotification) error
	ListAll(wsID uuid.UUID) ([]models.WorkspaceNotification, error)
	Delete(ID uuid.UUID) error
}
