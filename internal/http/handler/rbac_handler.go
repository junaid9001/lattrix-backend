package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

type RbacHandler struct {
	rbacService *services.RbacService
}

func NewRbacHandler(rbacService *services.RbacService) *RbacHandler {
	return &RbacHandler{rbacService: rbacService}
}

type createRoleRequestDTO struct {
	Name          string      `json:"name" validate:"required,min=2"`
	PermissionIDs []uuid.UUID `json:"permission_ids" validate:"required,min=1,dive,required"`
}

func (h *RbacHandler) CreateRoleAndAssignPermission(c *fiber.Ctx) error {

	var dto createRoleRequestDTO

	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	workspaceIDval := c.Locals("workspaceID")

	workspaceID, ok := workspaceIDval.(uuid.UUID)

	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "invalid workspaceID")
	}

	err := h.rbacService.CreateRoleAndAssignPermissions(workspaceID, dto.Name, dto.PermissionIDs)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Role created and permissions assigned successfully",
	})

}

func (h *RbacHandler) GetAllRoles(c *fiber.Ctx) error {
	workspaceIDval := c.Locals("workspaceID")
	workspaceID, ok := workspaceIDval.(uuid.UUID)

	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "invalid workspaceID")
	}

	roles, err := h.rbacService.AllRoles(workspaceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    roles,
	})

}

func (h *RbacHandler) GetAllPermissions(c *fiber.Ctx) error {
	permissions, err := h.rbacService.AllPermissions()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    permissions,
	})
}

func (h *RbacHandler) UpdateUserRole(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("userId"))
	workspaceID := c.Locals("workspaceID").(uuid.UUID)

	var body struct {
		RoleID uuid.UUID `json:"role_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(400, "invalid body")
	}

	err := h.rbacService.AssignRoleToUser(uint(userID), workspaceID, body.RoleID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(fiber.Map{"success": true})
}
