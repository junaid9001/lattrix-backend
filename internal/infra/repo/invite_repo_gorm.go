package repo

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type InvitationRepo struct {
	db *gorm.DB
}

func NewInvitationRepo(db *gorm.DB) repository.InvitationRepository {
	return &InvitationRepo{db: db}
}

func (r *InvitationRepo) WithTx(tx *gorm.DB) repository.InvitationRepository {
	return &InvitationRepo{db: tx}
}

func (r *InvitationRepo) Create(invite *models.WorkspaceInvitation) error {
	return r.db.Create(invite).Error
}

func (r *InvitationRepo) FindByID(id uuid.UUID) (*models.WorkspaceInvitation, error) {
	var invite models.WorkspaceInvitation
	err := r.db.First(&invite, "id = ?", id).Error
	return &invite, err
}

func (r *InvitationRepo) FindByToken(token string) (*models.WorkspaceInvitation, error) {
	var invite models.WorkspaceInvitation
	err := r.db.First(&invite, "token = ?", token).Error
	return &invite, err
}

func (r *InvitationRepo) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.
		Model(&models.WorkspaceInvitation{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}
