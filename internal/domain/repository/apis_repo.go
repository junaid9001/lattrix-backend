package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
)

type ApiRepository interface {
	Create(*models.API) error
	Delete(ID uuid.UUID, ApiGroupID uuid.UUID) error
	Update(ID uuid.UUID, ApiGroupID uuid.UUID, updates map[string]any) (*models.API, error)
	GetByID(ID uuid.UUID, ApiGroupID uuid.UUID) (*models.API, error)
	ListByGroup(ApiGroupID uuid.UUID) ([]models.API, error)
	ListActive() ([]models.API, error)
	UpdateStatus(ID uuid.UUID, ApiGroupID uuid.UUID, lastStatus string, lastCheckedAt time.Time) error
	ListDueForCheck(now time.Time) ([]models.API, error)
}
