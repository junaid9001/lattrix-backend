package repo

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
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

func (r *userRepository) UpdateProfile(userID int, updates map[string]interface{}) (*models.User, error) {

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
func (r *userRepository) CreateWorkSpace(userID uint) (uuid.UUID, error) {
	workSpace := models.Workspace{
		ID:      uuid.New(),
		OwnerID: userID,
	}
	err := r.db.Create(&workSpace).Error
	if err != nil {
		return uuid.Nil, err
	}

	return workSpace.ID, nil
}

func (r *userRepository) WithDB(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}
