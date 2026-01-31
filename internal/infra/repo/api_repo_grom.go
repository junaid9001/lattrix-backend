package repo

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type ApiRepo struct {
	db *gorm.DB
}

func NewApiRepo(db *gorm.DB) repository.ApiRepository {
	return &ApiRepo{db: db}
}

func (r *ApiRepo) Create(api *models.API) error {
	return r.db.Create(api).Error
}

func (r *ApiRepo) Delete(ID uuid.UUID, ApiGroupID uuid.UUID, workspaceID uuid.UUID) error {
	result := r.db.Where("id = ? AND api_group_id = ? AND workspace_id = ?", ID, ApiGroupID, workspaceID).Delete(&models.API{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("api not found or access denied")
	}

	return nil
}

func (r *ApiRepo) Update(ID uuid.UUID, ApiGroupID uuid.UUID, updates map[string]any, workspaceID uuid.UUID) (*models.API, error) {
	result := r.db.
		Model(&models.API{}).
		Where("id = ? AND api_group_id = ? AND workspace_id = ?", ID, ApiGroupID, workspaceID).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("api not found or access denied")
	}

	var api models.API

	if err := r.db.First(&api, ID).Error; err != nil {
		return nil, err
	}

	return &api, nil

}

func (r *ApiRepo) GetByID(ID uuid.UUID, ApiGroupID uuid.UUID) (*models.API, error) {
	var api models.API
	err := r.db.Where("id = ? AND api_group_id = ?", ID, ApiGroupID).First(&api).Error

	if err != nil {
		return nil, err
	}

	return &api, nil
}

func (r *ApiRepo) ListByGroup(ApiGroupID uuid.UUID) ([]models.API, error) {
	var apis []models.API

	err := r.db.Where("api_group_id = ?", ApiGroupID).Find(&apis).Error
	if err != nil {
		return nil, err
	}

	return apis, nil
}

func (r *ApiRepo) ListActive() ([]models.API, error) {
	var apis []models.API

	err := r.db.Where("is_active = ?", true).Find(&apis).Error
	if err != nil {
		return nil, err
	}

	return apis, nil
}

func (r *ApiRepo) UpdateStatus(ID uuid.UUID, ApiGroupID uuid.UUID, lastStatus string, lastCheckedAt time.Time) error {
	result := r.db.Model(&models.API{}).
		Where("id = ? AND api_group_id = ?", ID, ApiGroupID).
		Updates(map[string]any{
			"last_status":     lastStatus,
			"last_checked_at": lastCheckedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *ApiRepo) ListDueForCheck(now time.Time) ([]models.API, error) {
	var apis []models.API

	// result := r.db.Where("is_active = ?", true).
	// 	Where("last_checked_at IS NULL OR last_checked_at + (interval_seconds::text || ' seconds')::interval <= ?", now).
	// 	Find(&apis)

	result := r.db.Where("is_active = ?", true).Where("next_check_at <= ?", now).Order("next_check_at ASC").Find(&apis)

	if result.Error != nil {
		return nil, result.Error
	}

	return apis, nil
}

//plan

func (r *ApiRepo) CountByOwnerID(ownerID uint) (int64, error) {
	var count int64
	err := r.db.Table("apis").
		Joins("JOIN workspaces ON workspaces.id=apis.workspace_id").
		Where("workspaces.owner_id=?", ownerID).
		Count(&count).Error
	return count, err
}

// on exp downgrade
func (r *ApiRepo) EnforcePlanLimits(userID uint, maxApis int64, minInterval int) error {

	err := r.db.Exec(`
		UPDATE apis
		SET interval_seconds = ?
		FROM workspaces
		WHERE apis.workspace_id = workspaces.id
		AND workspaces.owner_id = ?
		AND apis.interval_seconds < ?
	`, minInterval, userID, minInterval).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ApiRepo) GetCheckHistory(apiID uuid.UUID, limit int) ([]models.ApiCheckResult, error) {
	var results []models.ApiCheckResult
	err := r.db.Where("api_id = ?", apiID).
		Order("checked_at desc").
		Limit(limit).
		Find(&results).Error
	return results, err
}
