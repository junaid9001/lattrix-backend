package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
)

type WorkspaceNotiService struct {
	WorkNotiRepo repository.WorkspaceNotificationRepository
}

func NewWorkspaceNotiService(workNotiRepo repository.WorkspaceNotificationRepository) *WorkspaceNotiService {
	return &WorkspaceNotiService{WorkNotiRepo: workNotiRepo}
}

func (s *WorkspaceNotiService) Create(wsID uuid.UUID, message string, title string) error {
	noti := &models.WorkspaceNotification{
		ID:          uuid.New(),
		WorkspaceID: wsID,
		Message:     message,
		Title:       title,
		CreatedAt:   time.Now(),
	}
	return s.WorkNotiRepo.Create(noti)
}

func (s *WorkspaceNotiService) ListALL(wsID uuid.UUID) ([]models.WorkspaceNotification, error) {
	notis, err := s.WorkNotiRepo.ListAll(wsID)
	if err != nil {
		return nil, err
	}
	return notis, nil
}
