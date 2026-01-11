package repository

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"gorm.io/gorm"
)

type RBACrepository interface {
	CreateRole(*models.Role) error

	AssignPermissionToRole(roleID uuid.UUID, permissionIDs []uuid.UUID) error

	AssignRoleToUser(userID uint, roleID, workspaceID uuid.UUID) error
	AllRoles(workspaceID uuid.UUID) (*[]models.Role, error)

	UserHasPermission(userID uint, workspaceID uuid.UUID, permissionCode string) (bool, error)
	UserPermissions(userID uint, workspaceID uuid.UUID) ([]string, error)

	PermissionsExist(permissionIDs []uuid.UUID) (bool, error)

	WithTx(tx *gorm.DB) RBACrepository
}
