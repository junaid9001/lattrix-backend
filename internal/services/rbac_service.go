package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type RbacService struct {
	RbacRepo repository.RBACrepository
	db       *gorm.DB
}

func NewRbacService(RbacRepo repository.RBACrepository, db *gorm.DB) *RbacService {
	return &RbacService{RbacRepo: RbacRepo, db: db}
}

func (s *RbacService) CreateRoleAndAssignPermissions(userID uint, workspaceID uuid.UUID, roleName string,
	permissionIDs []uuid.UUID) error {

	err := s.db.Transaction(func(tx *gorm.DB) error {

		rbacRepo := s.RbacRepo.WithTx(tx)

		ok, err := rbacRepo.PermissionsExist(permissionIDs)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("one or more permissions are invalid")
		}

		role := &models.Role{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			Name:        roleName,
		}
		if err := rbacRepo.CreateRole(role); err != nil {
			return err
		}

		if err := rbacRepo.AssignPermissionToRole(role.ID, permissionIDs); err != nil {
			return err
		}

		if err := rbacRepo.AssignRoleToUser(userID, role.ID, workspaceID); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *RbacService) AllRoles(workspaceID uuid.UUID) (*[]models.Role, error) {
	return s.RbacRepo.AllRoles(workspaceID)
}
