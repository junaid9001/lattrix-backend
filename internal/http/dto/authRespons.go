package dto

import "github.com/google/uuid"

type UserWorkspaceResponse struct {
	WorkspaceID uuid.UUID `json:"workspace_id"`
	Name        string    `json:"name"`
	Role        string    `json:"role"`
}

type LoginResponse struct {
	Success    bool                    `json:"success"`
	Workspaces []UserWorkspaceResponse `json:"workspaces"`
}

type SelectWorkspaceRequest struct {
	WorkspaceID uuid.UUID `json:"workspace_id" validate:"required"`
}
