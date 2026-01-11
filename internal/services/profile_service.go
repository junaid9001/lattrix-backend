package services

import (
	"errors"

	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
)

type ProfileService struct {
	userRepo repository.UserRepository
}

func NewProfileService(userRepo repository.UserRepository) *ProfileService {
	return &ProfileService{userRepo: userRepo}

}

func (s *ProfileService) GetUserProfile(userID int) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *ProfileService) UpdateProfileByID(userID int, username, email *string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if username != nil && user.Username != *username {
		updates["username"] = *username
	}

	if email != nil && user.Email != *email {
		updates["email"] = *email
	}

	if len(updates) == 0 {
		return user, errors.New("nothing to update")
	}

	user, err = s.userRepo.UpdateProfile(userID, updates)
	if err != nil {
		return nil, err
	}
	return user, nil

}
