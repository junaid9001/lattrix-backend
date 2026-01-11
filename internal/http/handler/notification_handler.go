package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
)

type NotificationHandler struct {
	repo repository.NotificationRepository
}

func NewNotificationHandler(repo repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

func (h *NotificationHandler) GetMyNotifications(c *fiber.Ctx) error {
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	list, err := h.repo.GetByUserID(uint(userID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch notifications")
	}

	return c.JSON(fiber.Map{"success": true, "data": list})
}
