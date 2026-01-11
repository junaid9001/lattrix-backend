package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
)

func RbacRoute(app *fiber.App, h *handler.RbacHandler) {
	rbac := app.Group("/rbac")
	rbac.Use(middleware.Auth())

	rbac.Post("/role", h.CreateRoleAndAssignPermission)
	rbac.Get("/roles", h.GetAllRoles)

}
