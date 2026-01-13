package repo

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type RbacRepo struct {
	db *gorm.DB
}

func NewRbacRepo(db *gorm.DB) repository.RBACrepository {
	return &RbacRepo{db: db}
}

func (r *RbacRepo) CreateRole(role *models.Role) error {
	return r.db.Create(&role).Error
}

func (r *RbacRepo) AssignPermissionToRole(roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	rolePermissions := make([]models.RolePermission, 0, len(permissionIDs))

	for _, val := range permissionIDs {
		rolePermissions = append(rolePermissions, models.RolePermission{
			ID:           uuid.New(),
			RoleID:       roleID,
			PermissionID: val,
		})
	}
	return r.db.Create(&rolePermissions).Error
}

func (r *RbacRepo) AssignRoleToUser(userID uint, roleID, workspaceID uuid.UUID) error {
	var role models.Role
	if err := r.db.Where("id = ?", roleID).First(&role).Error; err != nil {
		return err
	}
	var userRole models.UserRole

	err := r.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&userRole).Error

	if err == nil {
		userRole.RoleID = roleID
		return r.db.Save(&userRole).Error
	} else {
		newUserRole := models.UserRole{
			ID:          uuid.New(),
			UserID:      userID,
			RoleID:      roleID,
			WorkspaceID: workspaceID,
		}
		if err := r.db.Create(&newUserRole).Error; err != nil {
			return err
		}

	}
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role.Name).Error
}

func (r *RbacRepo) AllRoles(workspaceID uuid.UUID) ([]models.Role, error) {
	var roles []models.Role

	err := r.db.Where("workspace_id = ?", workspaceID).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	r.db.Debug()
	return roles, nil
}

func (r *RbacRepo) UserHasPermission(
	userID uint,
	workspaceID uuid.UUID,
	permissionCode string,
) (bool, error) {

	var count int64

	err := r.db.
		Table("user_roles ur").
		Joins("JOIN roles r ON r.id = ur.role_id AND r.deleted_at IS NULL").
		Joins("JOIN role_permissions rp ON rp.role_id = r.id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("ur.user_id = ?", userID).
		Where("ur.workspace_id = ?", workspaceID).
		Where("p.code = ?", permissionCode).
		Limit(1).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *RbacRepo) UserPermissions(userID uint, workspaceID uuid.UUID) ([]string, error) {
	var permissions []string

	err := r.db.
		Table("user_roles ur").
		Select("DISTINCT p.code").
		Joins("JOIN roles r ON r.id = ur.role_id AND r.deleted_at IS NULL").
		Joins("JOIN role_permissions rp ON rp.role_id = r.id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("ur.user_id = ?", userID).
		Where("ur.workspace_id = ?", workspaceID).
		Pluck("p.code", &permissions).
		Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *RbacRepo) PermissionsExist(permissionIDs []uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Table("permissions").
		Where("id IN ?", permissionIDs).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count == int64(len(permissionIDs)), nil
}

func (r *RbacRepo) WithTx(tx *gorm.DB) repository.RBACrepository {
	return &RbacRepo{db: tx}
}

func (r *RbacRepo) AllPermissions() ([]models.Permission, error) {
	var permissions []models.Permission

	err := r.db.Find(&permissions).Error
	return permissions, err
}
