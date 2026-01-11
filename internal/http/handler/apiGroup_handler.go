package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

type ApiGroupHandler struct {
	apiGroupService *services.ApiGroupService
}

func NewApiGroupHandler(apiGroupService *services.ApiGroupService) *ApiGroupHandler {
	return &ApiGroupHandler{apiGroupService: apiGroupService}
}

type createGroupdto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (h *ApiGroupHandler) CreateNewApiGroupHandler(c *fiber.Ctx) error {
	var req createGroupdto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid JSON"})
	}
	val, val2 := c.Locals("workspaceID"), c.Locals("userID")
	workspaceID, ok := val.(uuid.UUID)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized: invalid workspace context",
		})
	}

	userID, ok := val2.(int)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized: invalid userid context",
		})
	}

	apigrp, err := h.apiGroupService.CreateNewApiGroup(userID, req.Name, req.Description, workspaceID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": apigrp})
}

type deleteGroupdto struct {
	ID uuid.UUID `json:"api_group_id" validate:"required"`
}

func (h *ApiGroupHandler) DeleteApiGroupHandler(c *fiber.Ctx) error {

	var req deleteGroupdto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false})
	}
	val := c.Locals("workspaceID")
	workspaceID, ok := val.(uuid.UUID)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized: invalid workspace context",
		})
	}

	err := h.apiGroupService.DeleteApiGroup(req.ID, workspaceID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "API group deleted successfully",
	})
}

func (h *ApiGroupHandler) GetApiGroupHandler(c *fiber.Ctx) error {

	idParam := c.Params("id")
	groupID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid uuid format",
		})
	}

	val := c.Locals("workspaceID")
	workspaceID, ok := val.(uuid.UUID)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized: invalid workspace context",
		})
	}

	apigrp, err := h.apiGroupService.Getapigroupbyid(groupID, workspaceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "api group not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    apigrp,
	})

}

type updateGropdto struct {
	ID          uuid.UUID `json:"api_group_id" validate:"required"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
}

func (h *ApiGroupHandler) UpdateApiGroupHandler(c *fiber.Ctx) error {
	var req updateGropdto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false})
	}
	val := c.Locals("workspaceID")
	workspaceID, ok := val.(uuid.UUID)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized: invalid workspace context",
		})
	}

	apigrp, err := h.apiGroupService.Updateapigroup(req.ID, workspaceID, req.Name, req.Description)
	if err != nil {
		if err.Error() == "nothing to update" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to update api group",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "API group updated successfully",
		"data":    apigrp,
	})

}
