package repo

import (
	"time"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type NotificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) repository.NotificationRepository {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) WithTx(tx *gorm.DB) repository.NotificationRepository {
	return &NotificationRepo{db: tx}
}

func (r *NotificationRepo) Create(n *models.Notification) error {
	return r.db.Create(n).Error
}

func (r *NotificationRepo) GetByUserID(userID uint) ([]models.Notification, error) {
	var list []models.Notification
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&list).
		Error
	return list, err
}

func (r *NotificationRepo) MarkAsReadByReference(refID uuid.UUID, notifType string) error {
	now := time.Now()
	return r.db.
		Model(&models.Notification{}).
		Where("reference_id = ? AND type = ?", refID, notifType).
		Update("read_at", &now).
		Error
}
