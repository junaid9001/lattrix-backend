package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(service *services.AdminService) *AdminHandler {
	return &AdminHandler{adminService: service}
}

func (h *AdminHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.adminService.GetDashboardStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch stats", "details": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": stats})
}

func (h *AdminHandler) GetActivities(c *fiber.Ctx) error {
	activities, err := h.adminService.GetRecentStripeActivities()
	if err != nil {

		return c.JSON(fiber.Map{"success": false, "data": []interface{}{}, "error": "Stripe Sync Failed"})
	}
	return c.JSON(fiber.Map{"success": true, "data": activities})
}

func (h *AdminHandler) GetSystemHealth(c *fiber.Ctx) error {
	health := h.adminService.GetSystemHealth()
	return c.JSON(fiber.Map{"success": true, "data": health})
}

func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, total, err := h.adminService.GetAllUsers(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *AdminHandler) ToggleBan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.adminService.ToggleUserBan(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to update user status"})
	}

	status := "Active"
	if !user.IsActive {
		status = "Banned"
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   "User is now " + status,
		"is_active": user.IsActive,
	})
}
