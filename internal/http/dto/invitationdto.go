package dto

import "github.com/google/uuid"

type AcceptInvitationRequestDTO struct {
	Token string `json:"token"`
}

type SendInvitationRequestDTO struct {
	Email  string    `json:"email" validate:"required,email"`
	RoleID uuid.UUID `json:"role_id" validate:"required"`
}
