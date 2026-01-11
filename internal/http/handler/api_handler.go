package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/http/dto"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

type ApiHandler struct {
	apiService *services.ApiService
}

func NewApiHandler(apiService *services.ApiService) *ApiHandler {
	return &ApiHandler{apiService: apiService}
}

func (h *ApiHandler) RegisterHandler(c *fiber.Ctx) error {

	var dto dto.ApiRegisterDto

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "invalid body"})
	}

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: invalid user context",
		})
	}

	workspaceIDVal := c.Locals("workspaceID")
	workspaceID, ok := workspaceIDVal.(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: invalid workspace context",
		})
	}

	apiGroupIDParam := c.Params("api-group-id")
	apiGroupID, err := uuid.Parse(apiGroupIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid API group ID",
		})
	}

	api, err := h.apiService.RegisterApiService(uint(userID), apiGroupID, workspaceID, &dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "data": api})

}

// update api
func (h *ApiHandler) UpdateApi(c *fiber.Ctx) error {
	var dto dto.ApiUpdateDto

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "invalid body"})
	}

	ApiIDParam := c.Params("api-id")
	ApiGroupIDParam := c.Params("api-group-id")

	apiID, err := uuid.Parse(ApiIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid API ID",
		})
	}

	apiGroupID, err := uuid.Parse(ApiGroupIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid API group ID",
		})
	}

	api, err := h.apiService.UpdateApi(apiID, apiGroupID, dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "data": api})

}

// delte api
func (h *ApiHandler) Delete(c *fiber.Ctx) error {
	apiID, err := uuid.Parse(c.Params("api-id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid API ID"})
	}

	groupID, err := uuid.Parse(c.Params("api-group-id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid API group ID"})
	}

	if err := h.apiService.DeleteApi(apiID, groupID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete API",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "API deleted successfully"})
}

//list by group

func (h *ApiHandler) ListByGroup(c *fiber.Ctx) error {

	apiGroupIDParam := c.Params("api-group-id")
	apiGroupID, err := uuid.Parse(apiGroupIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid API group ID",
		})
	}

	apis, err := h.apiService.ListApisByGroup(apiGroupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    apis,
	})
}
