package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"github.com/junaid9001/lattrix-backend/internal/utils/jwtutil"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo     repository.UserRepository
	apiGroupRepo repository.ApiGroupRepository
	db           *gorm.DB
}

func NewAuthSevice(userRepo repository.UserRepository, apiGroupRepo repository.ApiGroupRepository, db *gorm.DB) *AuthService {
	return &AuthService{userRepo: userRepo, apiGroupRepo: apiGroupRepo, db: db}
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

		user := &models.User{
			Username: username,
			Email:    email,
			Password: string(hashed),
			Role:     "superadmin",
		}

		if err := userRepo.Create(user); err != nil {
			return err
		}
		uuidd, err := userRepo.CreateWorkSpace(user.ID)
		if err != nil {
			return err
		}
		user.WorkspaceID = uuidd
		tx.Save(&user)

		mainApiGroup := &models.ApiGroup{
			ID:             uuid.New(),
			WorkspaceID:    user.WorkspaceID,
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

func (s *AuthService) Login(email, password string) (string, string, error) {

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid email or password")
	}

	access, err := jwtutil.CreateAccessToken(int(user.ID), user.WorkspaceID.String(), user.Role)
	if err != nil {
		return "", "", err
	}
	refresh, err := jwtutil.CreateRefreshToken(int(user.ID))
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil

}

func (s *AuthService) RefreshAccessToken(id int) (string, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !user.IsActive {
		return "", errors.New("user account is deactivated")
	}

	return jwtutil.CreateAccessToken(int(user.ID), user.WorkspaceID.String(), user.Role)
}
