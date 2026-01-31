package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

func ApiRoutes(app *fiber.App, h *handler.ApiHandler, rbacService *services.RbacService) {

	// group-level routes (need api-group-id)
	apiGroup := app.Group(
		"/api-groups/:apiGroupID",
	)

	apiGroup.Use(middleware.Auth())

	//create api
	apiGroup.Post("/apis", middleware.RequirePermission(rbacService, "api:create"), h.RegisterHandler)

	//list apis by group
	apiGroup.Get("/apis", h.ListByGroup)

	//update api
	apiGroup.Put("/apis/:apiId", middleware.RequirePermission(rbacService, "api:update"), h.UpdateApi)

	//delete api
	apiGroup.Delete("/apis/:apiId", middleware.RequirePermission(rbacService, "api:delete"), h.Delete)

	apiGroup.Get("/apis/:api_id/history", middleware.RequirePermission(rbacService, "api:read"), h.GetMetricsHistory)
}
