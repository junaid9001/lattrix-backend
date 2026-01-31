package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"gorm.io/gorm"
)

func RequireSuperAdmin(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID")

		var user models.User

		if err := db.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		if !user.IsSuperAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access Denied: Admins only"})
		}

		return c.Next()
	}
}
