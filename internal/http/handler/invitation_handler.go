package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/http/dto"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

type InvitationHandler struct {
	services *services.InvitationService
}

func NewInvitationHandler(service *services.InvitationService) *InvitationHandler {
	return &InvitationHandler{services: service}
}

func (h *InvitationHandler) SendInvite(c *fiber.Ctx) error {
	var req dto.SendInvitationRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	workspaceIDVal := c.Locals("workspaceID")
	workspaceID, ok := workspaceIDVal.(uuid.UUID)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized: invalid workspace context")
	}

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized: invalid user context")
	}

	err := h.services.CreateNewInvite(workspaceID, req.RoleID, req.Email, uint(userID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "invitation sent",
	})
}

func (h *InvitationHandler) AcceptInvitation(c *fiber.Ctx) error {

	var dto dto.AcceptInvitationRequestDTO
	if err := c.BodyParser(&dto); err != nil || dto.Token == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	email := c.Locals("email").(string)

	jwt, err := h.services.AcceptInvitation(uint(userID), email, dto.Token)
	if err != nil {
		switch err {
		case services.ErrInvalidInvite, services.ErrInviteExpired:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case services.ErrInviteProcessed:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "internal error")
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    jwt,
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{
		"message": "invitation accepted",
	})
}
