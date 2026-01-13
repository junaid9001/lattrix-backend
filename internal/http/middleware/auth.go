package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/utils/jwtutil"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("access_token")

		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"success": false, "error": "missing access token"})
		}

		claims, err := jwtutil.ValidateAccessToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"success": false, "error": "invalid token"})
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"success": false, "error": "invalid subject"})
		}

		workspaceID, _ := uuid.Parse(claims.WorkspaceID)

		c.Locals("userID", userID)
		c.Locals("role", claims.Role)
		c.Locals("workspaceID", workspaceID)

		return c.Next()
	}
}
