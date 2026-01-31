package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"github.com/junaid9001/lattrix-backend/internal/http/dto"
	"github.com/junaid9001/lattrix-backend/internal/utils/jwtutil"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo     repository.UserRepository
	apiGroupRepo repository.ApiGroupRepository
	rbacRepo     repository.RBACrepository
	db           *gorm.DB
}

func NewAuthSevice(userRepo repository.UserRepository, apiGroupRepo repository.ApiGroupRepository, rbacRepo repository.RBACrepository, db *gorm.DB) *AuthService {
	return &AuthService{userRepo: userRepo, apiGroupRepo: apiGroupRepo, rbacRepo: rbacRepo, db: db}
}

func (s *AuthService) SignUP(username, email, password string) error {
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		userRepo := s.userRepo.WithDB(tx)
		apiGroupRepo := s.apiGroupRepo.WithDB(tx)
		rbacRepo := s.rbacRepo.WithTx(tx)

		user := &models.User{
			Username: username,
			Email:    email,
			Password: string(hashed),
		}

		if err := userRepo.Create(user); err != nil {
			return err
		}

		wsName := username + "'s Workspace"

		wsID, err := userRepo.CreateWorkSpace(user.ID, wsName)
		if err != nil {
			return err
		}

		ownerRole := &models.Role{
			ID:          uuid.New(),
			WorkspaceID: wsID,
			Name:        "Owner",
		}
		if err := rbacRepo.CreateRole(ownerRole); err != nil {
			return err
		}

		var superPerm models.Permission
		if err := tx.Where("code = ?", "role:superadmin").First(&superPerm).Error; err != nil {
			return errors.New("system error: superadmin permission not seeded")
		}

		if err := rbacRepo.AssignPermissionToRole(ownerRole.ID, []uuid.UUID{superPerm.ID}); err != nil {
			return err
		}

		if err := rbacRepo.AssignRoleToUser(user.ID, ownerRole.ID, wsID); err != nil {
			return err
		}

		mainApiGroup := &models.ApiGroup{
			ID:             uuid.New(),
			WorkspaceID:    wsID,
			Name:           "main",
			CreatedByID:    user.ID,
			CreatedByEmail: user.Email,
			Description:    "Default Group",
		}

		if err := apiGroupRepo.Create(mainApiGroup); err != nil {
			return err
		}

		return nil
	})
}

func (s *AuthService) Login(email, password string) (string, []dto.UserWorkspaceResponse, error) {

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	refresh, err := jwtutil.CreateRefreshToken(int(user.ID))
	if err != nil {
		return "", nil, err
	}

	workspace, err := s.userRepo.GetUserWorkspaces(user.ID)
	if err != nil {
		return "", nil, err
	}

	return refresh, workspace, nil

}

// gen accesstoken for workspace
func (s *AuthService) GenerateAccessTokenForWorkspace(userID int, workspaceID uuid.UUID) (string, error) {
	var userRole models.UserRole
	if err := s.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&userRole).Error; err != nil {
		return "", errors.New("access denied to this workspace")
	}
	var role models.Role
	if err := s.db.First(&role, userRole.RoleID).Error; err != nil {
		return "", errors.New("role not found")
	}

	return jwtutil.CreateAccessToken(userID, workspaceID.String(), role.Name)
}

func (s *AuthService) RefreshAccessToken(id int) (string, error) {
	// user, err := s.userRepo.FindByID(id)
	// if err != nil {
	// 	return "", errors.New("user not found")
	// }

	// if !user.IsActive {
	// 	return "", errors.New("user account is deactivated")
	// }

	return "", errors.New("use select-workspace endpoint")
}

func (s *AuthService) CreateWorkspace(userID uint, name string) (uuid.UUID, error) {
	var workspaceID uuid.UUID

	err := s.db.Transaction(func(tx *gorm.DB) error {
		userRepo := s.userRepo.WithDB(tx)
		rbacRepo := s.rbacRepo.WithTx(tx)
		apiGroupRepo := s.apiGroupRepo.WithDB(tx)

		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		id, err := userRepo.CreateWorkSpace(userID, name)
		if err != nil {
			return err
		}
		workspaceID = id

		ownerRole := &models.Role{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			Name:        "Owner",
		}
		if err := rbacRepo.CreateRole(ownerRole); err != nil {
			return err
		}

		var superPerm models.Permission
		if err := tx.Where("code = ?", "role:superadmin").First(&superPerm).Error; err != nil {
			return errors.New("system error: superadmin permission not seeded")
		}

		if err := rbacRepo.AssignPermissionToRole(ownerRole.ID, []uuid.UUID{superPerm.ID}); err != nil {
			return err
		}

		if err := rbacRepo.AssignRoleToUser(userID, ownerRole.ID, workspaceID); err != nil {
			return err
		}

		mainApiGroup := &models.ApiGroup{
			ID:             uuid.New(),
			WorkspaceID:    workspaceID,
			Name:           "main",
			CreatedByID:    user.ID,
			CreatedByEmail: user.Email,
			Description:    "Default Group",
		}

		if err := apiGroupRepo.Create(mainApiGroup); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	return workspaceID, nil
}

func (s *AuthService) UserWorkspaces(userID uint) ([]dto.UserWorkspaceResponse, error) {
	workSpaces, err := s.userRepo.GetUserWorkspaces(userID)
	if err != nil {
		return nil, err
	}
	return workSpaces, nil
}
