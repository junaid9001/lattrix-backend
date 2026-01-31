package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
)

type NotificationHandler struct {
	repo              repository.NotificationRepository
	workspaceNotiRepo repository.WorkspaceNotificationRepository
}

func NewNotificationHandler(repo repository.NotificationRepository, workspaceNotiRepo repository.WorkspaceNotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo, workspaceNotiRepo: workspaceNotiRepo}
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

func (h *NotificationHandler) WorkspaceNotifications(c *fiber.Ctx) error {
	wsID := c.Locals("workspaceID").(uuid.UUID)
	notis, err := h.workspaceNotiRepo.ListAll(wsID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    notis,
	})
}

func (h *NotificationHandler) DeleteWorkspaceNotification(c *fiber.Ctx) error {
	notificationIDParam := c.Params("notificationId")
	notificationID, err := uuid.Parse(notificationIDParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid notification ID")
	}

	if err := h.workspaceNotiRepo.Delete(notificationID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete notification")
	}

	return c.JSON(fiber.Map{"success": true, "message": "Notification deleted successfully"})

}
