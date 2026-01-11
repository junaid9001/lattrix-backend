package repository

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	WithTx(tx *gorm.DB) NotificationRepository

	Create(n *models.Notification) error
	GetByUserID(userID uint) ([]models.Notification, error)
	MarkAsReadByReference(refID uuid.UUID, notifType string) error
}
