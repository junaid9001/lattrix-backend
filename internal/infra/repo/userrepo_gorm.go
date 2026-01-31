package repo

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"github.com/junaid9001/lattrix-backend/internal/http/dto"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email=?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ID int) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, ID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateProfile(userID uint, updates map[string]interface{}) (*models.User, error) {

	if err := r.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return nil, err
	}

	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// only during signup oneuser/oneworkspace for now
func (r *userRepository) CreateWorkSpace(userID uint, name string) (uuid.UUID, error) {
	workSpace := models.Workspace{
		ID:      uuid.New(),
		Name:    name,
		OwnerID: userID,
	}
	err := r.db.Create(&workSpace).Error
	if err != nil {
		return uuid.Nil, err
	}

	return workSpace.ID, nil
}

func (r *userRepository) WorkspaceUsers(
	workspaceID uuid.UUID,
) ([]dto.WorkspaceUsers, error) {

	var users []dto.WorkspaceUsers

	err := r.db.
		Table("users u").
		Select(`
			u.id AS user_id,
			u.email,
			r.id AS role_id,
			r.name AS role
		`).
		Joins("JOIN user_roles ur ON ur.user_id = u.id").
		Joins("JOIN roles r ON r.id = ur.role_id AND r.deleted_at IS NULL").
		Where("ur.workspace_id = ?", workspaceID).
		Scan(&users).
		Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// get a users all workspace
func (r *userRepository) GetUserWorkspaces(userID uint) ([]dto.UserWorkspaceResponse, error) {
	var results []dto.UserWorkspaceResponse

	err := r.db.Table("workspaces w").
		Select("w.id as workspace_id, w.name, r.name as role").
		Joins("JOIN user_roles ur ON ur.workspace_id = w.id").
		Joins("JOIN roles r ON r.id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Scan(&results).Error

	return results, err
}

func (r *userRepository) WithDB(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}
