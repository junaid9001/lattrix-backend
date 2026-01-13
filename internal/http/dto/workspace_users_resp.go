package dto

import "github.com/google/uuid"

type WorkspaceUsers struct {
	UserID uint      `json:"user_id"`
	Email  string    `json:"email"`
	RoleID uuid.UUID `json:"role_id"`
	Role   string    `json:"role"`
}
