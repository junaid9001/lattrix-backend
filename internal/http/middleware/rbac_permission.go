package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

func RequirePermission(rbacService *services.RbacService, code string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		roleVal := c.Locals("role")

		role, ok := roleVal.(string)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		if role == "superadmin" {
			return c.Next()
		}

		userIDVal := c.Locals("userID")
		userIDInt, ok := userIDVal.(int)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		userID := uint(userIDInt)

		workspaceIDVal := c.Locals("workspaceID")
		workspaceID, ok := workspaceIDVal.(uuid.UUID)
		if !ok {
			return fiber.NewError(fiber.StatusBadRequest, "workspace not selected")
		}

		ok, err := rbacService.RbacRepo.UserHasPermission(userID, workspaceID, code)
		if err != nil {
			return fiber.NewError(
				fiber.StatusInternalServerError,
				"permission check failed",
			)
		}

		if !ok {
			return fiber.NewError(
				fiber.StatusForbidden,
				"insufficient permissions",
			)
		}

		return c.Next()

	}
}
