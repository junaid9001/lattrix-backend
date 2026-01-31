package repo

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type WorkspaceNotificationRepo struct {
	db *gorm.DB
}

func NewWorkspaceNotificationRepository(db *gorm.DB) repository.WorkspaceNotificationRepository {
	return &WorkspaceNotificationRepo{db: db}
}

func (r *WorkspaceNotificationRepo) Create(notification *models.WorkspaceNotification) error {
	return r.db.Create(&notification).Error
}

func (r *WorkspaceNotificationRepo) ListAll(wsID uuid.UUID) ([]models.WorkspaceNotification, error) {
	var notifications []models.WorkspaceNotification
	err := r.db.Where("workspace_id = ?", wsID).Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *WorkspaceNotificationRepo) Delete(ID uuid.UUID) error {
	return r.db.Delete(&models.WorkspaceNotification{}, ID).Error
}
