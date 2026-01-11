package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/dto"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

type UserProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileSevice *services.ProfileService) *UserProfileHandler {
	return &UserProfileHandler{profileService: profileSevice}
}

func (h *UserProfileHandler) GetProfile(c *fiber.Ctx) error {
	val := c.Locals("userID")
	if val == nil {
		return c.Status(401).JSON(fiber.Map{"success": false})
	}

	userID, ok := val.(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "invalid user context",
		})
	}

	user, err := h.profileService.GetUserProfile(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "user not found"})
	}

	return c.JSON(fiber.Map{"success": true, "data": fiber.Map{
		"username": user.Username,
		"email":    user.Email,
	}})
}

func (h *UserProfileHandler) UpdateProfileByID(c *fiber.Ctx) error {
	val := c.Locals("userID")
	if val == nil {
		return c.Status(401).JSON(fiber.Map{"success": false})
	}

	userID, ok := val.(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "invalid user context",
		})
	}

	var req dto.ProfileUpdateRequet

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "validation failed", "error": err.Error()})
	}

	user, err := h.profileService.UpdateProfileByID(userID, req.Username, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})

}
