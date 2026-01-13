package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

func RbacRoute(app *fiber.App, h *handler.RbacHandler, rbacService *services.RbacService) {
	rbac := app.Group("/rbac")
	rbac.Use(middleware.Auth())

	rbac.Post("/roles", middleware.RequirePermission(rbacService, "role:superadmin"), h.CreateRoleAndAssignPermission)
	rbac.Get("/roles", middleware.RequirePermission(rbacService, "role:superadmin"), h.GetAllRoles)

	rbac.Get("/permissions", middleware.RequirePermission(rbacService, "role:superadmin"), h.GetAllPermissions)
	rbac.Put("/users/:userId/role", middleware.RequirePermission(rbacService, "role:superadmin"), h.UpdateUserRole)

}
