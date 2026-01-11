package services

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"github.com/junaid9001/lattrix-backend/internal/http/dto"
)

type ApiService struct {
	apiRepo repository.ApiRepository
}

func NewApiService(apiRepo repository.ApiRepository) *ApiService {
	return &ApiService{apiRepo: apiRepo}
}

func (s *ApiService) RegisterApiService(
	userID uint,
	apiGroupID uuid.UUID,
	workspaceID uuid.UUID,
	dto *dto.ApiRegisterDto,
) (*models.API, error) {

	headers, err := json.Marshal(dto.Headers)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(dto.Body)
	if err != nil {
		return nil, err
	}

	expectedStatusCodes, err := json.Marshal(dto.ExpectedStatusCodes)
	if err != nil {
		return nil, err
	}

	interval := 60
	if dto.IntervalSeconds != nil {
		interval = *dto.IntervalSeconds
	}

	timeout := 3000
	if dto.TimeoutMs != nil {
		timeout = *dto.TimeoutMs
	}

	if dto.AuthType != "NONE" && dto.AuthValue == nil {
		return nil, errors.New("auth_value required")
	}

	id := uuid.New()

	api := models.API{
		ID:                     id,
		UserID:                 userID,
		ApiGroupID:             apiGroupID,
		WorkspaceID:            workspaceID,
		Name:                   dto.Name,
		Description:            dto.Description,
		URL:                    dto.URL,
		Method:                 dto.Method,
		AuthType:               dto.AuthType,
		AuthIn:                 dto.AuthIn,
		AuthKey:                dto.AuthKey,
		AuthValue:              dto.AuthValue,
		Headers:                headers,
		BodyType:               dto.BodyType,
		Body:                   body,
		IntervalSeconds:        interval,
		TimeoutMs:              timeout,
		IsActive:               true,
		ExpectedStatusCodes:    expectedStatusCodes,
		ExpectedResponseTimeMs: dto.ExpectedResponseTimeMs,
		ExpectedBodyContains:   dto.ExpectedBodyContains,
	}

	if err := s.apiRepo.Create(&api); err != nil {
		return nil, err
	}

	return &api, nil
}

func (s *ApiService) UpdateApi(
	ID uuid.UUID,
	apiGroupID uuid.UUID,
	dto dto.ApiUpdateDto,
) (*models.API, error) {

	api, err := s.apiRepo.GetByID(ID, apiGroupID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)

	if dto.Name != nil && api.Name != *dto.Name {
		updates["name"] = *dto.Name
	}

	if dto.Description != nil {
		updates["description"] = dto.Description
	}

	if dto.URL != nil && api.URL != *dto.URL {
		updates["url"] = *dto.URL
	}

	if dto.Method != nil && api.Method != *dto.Method {
		updates["method"] = *dto.Method
	}

	if dto.AuthType != nil && api.AuthType != *dto.AuthType {
		updates["auth_type"] = *dto.AuthType
	}

	if dto.AuthIn != nil {
		updates["auth_in"] = dto.AuthIn
	}

	if dto.AuthKey != nil {
		updates["auth_key"] = dto.AuthKey
	}

	if dto.AuthValue != nil {
		updates["auth_value"] = dto.AuthValue
	}

	if dto.Headers != nil {
		headersJSON, err := json.Marshal(*dto.Headers)
		if err != nil {
			return nil, err
		}
		updates["headers"] = headersJSON
	}

	if dto.BodyType != nil {
		updates["body_type"] = dto.BodyType
	}

	if dto.Body != nil {
		bodyJSON, err := json.Marshal(*dto.Body)
		if err != nil {
			return nil, err
		}
		updates["body"] = bodyJSON
	}

	if dto.IntervalSeconds != nil && api.IntervalSeconds != *dto.IntervalSeconds {
		updates["interval_seconds"] = *dto.IntervalSeconds
	}

	if dto.TimeoutMs != nil && api.TimeoutMs != *dto.TimeoutMs {
		updates["timeout_ms"] = *dto.TimeoutMs
	}

	if dto.ExpectedStatusCodes != nil {
		codesJSON, err := json.Marshal(*dto.ExpectedStatusCodes)
		if err != nil {
			return nil, err
		}
		updates["expected_status_codes"] = codesJSON
	}

	if dto.ExpectedResponseTimeMs != nil {
		updates["expected_response_time_ms"] = dto.ExpectedResponseTimeMs
	}

	if dto.ExpectedBodyContains != nil {
		updates["expected_body_contains"] = dto.ExpectedBodyContains
	}

	if len(updates) == 0 {
		return api, nil
	}

	updatedApi, err := s.apiRepo.Update(ID, apiGroupID, updates)
	if err != nil {
		return nil, err
	}

	return updatedApi, nil
}

func (s *ApiService) DeleteApi(
	ID uuid.UUID,
	apiGroupID uuid.UUID,
) error {

	_, err := s.apiRepo.GetByID(ID, apiGroupID)
	if err != nil {
		return err
	}

	if err := s.apiRepo.Delete(ID, apiGroupID); err != nil {
		return err
	}

	return nil
}

//list by groupid

func (s *ApiService) ListApisByGroup(apiGroupID uuid.UUID) ([]models.API, error) {
	return s.apiRepo.ListByGroup(apiGroupID)
}
