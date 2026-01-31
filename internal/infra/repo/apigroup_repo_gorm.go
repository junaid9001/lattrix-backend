package repo

import (
	"errors"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type apiGroupRepo struct {
	db *gorm.DB
}

func NewApiGroupRepository(db *gorm.DB) repository.ApiGroupRepository {
	return &apiGroupRepo{db: db}
}

func (r *apiGroupRepo) Create(apiGroup *models.ApiGroup) error {
	return r.db.Create(apiGroup).Error
}

func (r *apiGroupRepo) Delete(ID, workspaceID uuid.UUID) error {
	result := r.db.Where("id = ? AND workspace_id = ?", ID, workspaceID).Delete(&models.ApiGroup{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("api group not found or access denied")
	}
	return nil
}

func (r *apiGroupRepo) FindByID(ID, workspaceID uuid.UUID) (*models.ApiGroup, error) {
	var apiGroup models.ApiGroup

	err := r.db.
		Where("id = ? AND workspace_id = ?", ID, workspaceID).
		First(&apiGroup).Error

	if err != nil {
		return nil, err
	}

	return &apiGroup, nil
}

func (r *apiGroupRepo) Update(id uuid.UUID, workspaceID uuid.UUID, updates map[string]any) (*models.ApiGroup, error) {

	result := r.db.
		Model(&models.ApiGroup{}).
		Where("id = ? AND workspace_id = ?", id, workspaceID).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("api group not found or access denied")
	}

	var apiGroup models.ApiGroup
	r.db.First(&apiGroup, id) //not multitenat now

	return &apiGroup, nil
}

func (r *apiGroupRepo) ListGroups(workspaceID uuid.UUID) (*[]models.ApiGroup, error) {
	var apiGroups []models.ApiGroup
	err := r.db.Where("workspace_id = ?", workspaceID).Find(&apiGroups).Error
	if err != nil {
		return nil, err
	}
	return &apiGroups, nil
}

func (r *apiGroupRepo) WithDB(db *gorm.DB) repository.ApiGroupRepository {
	return &apiGroupRepo{db: db}
}
