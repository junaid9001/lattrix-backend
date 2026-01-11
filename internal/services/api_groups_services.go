package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
)

type ApiGroupService struct {
	apiGroupRepo repository.ApiGroupRepository
	userRepo     repository.UserRepository
}

func NewApiGroupService(apiGroupRepo repository.ApiGroupRepository, userRepo repository.UserRepository) *ApiGroupService {
	return &ApiGroupService{apiGroupRepo: apiGroupRepo, userRepo: userRepo}
}

func (s *ApiGroupService) CreateNewApiGroup(userID int, name, description string, workspaceID uuid.UUID) (*models.ApiGroup, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	apigrp := &models.ApiGroup{
		ID:             uuid.New(),
		WorkspaceID:    workspaceID,
		Name:           name,
		CreatedByEmail: user.Email,
		CreatedByID:    user.ID,
		Description:    description,
	}

	if user.WorkspaceID != workspaceID {
		return nil, errors.New("workspace mismatch")
	}

	if err := s.apiGroupRepo.Create(apigrp); err != nil {

		return nil, err
	}

	return apigrp, nil
}

func (s *ApiGroupService) DeleteApiGroup(ID, workspaceID uuid.UUID) error {
	err := s.apiGroupRepo.Delete(ID, workspaceID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ApiGroupService) Getapigroupbyid(ID, workspaceID uuid.UUID) (*models.ApiGroup, error) {

	apigrp, err := s.apiGroupRepo.FindByID(ID, workspaceID)

	if err != nil {
		return nil, err
	}

	return apigrp, nil
}

func (s *ApiGroupService) Updateapigroup(ID, workspaceID uuid.UUID, name, description *string) (*models.ApiGroup, error) {
	apigrp, err := s.apiGroupRepo.FindByID(ID, workspaceID)
	if err != nil {
		return nil, err
	}
	updates := make(map[string]any)

	if name != nil && apigrp.Name != *name {
		updates["name"] = *name
	}

	if description != nil && apigrp.Description != *description {
		updates["description"] = *description
	}

	if len(updates) == 0 {
		return nil, errors.New("nothing to update")
	}

	apigrp, err = s.apiGroupRepo.Update(ID, workspaceID, updates)

	if err != nil {
		return nil, err
	}
	return apigrp, nil
}
