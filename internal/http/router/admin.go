package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
	"gorm.io/gorm"
)

func AdminRoutes(app *fiber.App, h *handler.AdminHandler, db *gorm.DB) {
	// Protected Group: Auth + SuperAdmin Check
	admin := app.Group("/admin", middleware.Auth(), middleware.RequireSuperAdmin(db))

	admin.Get("/stats", h.GetStats)
	admin.Get("/activities", h.GetActivities)
	admin.Get("/health", h.GetSystemHealth)

	admin.Get("/users", h.ListUsers)
	admin.Patch("/users/:id/ban", h.ToggleBan)
}
