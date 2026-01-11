// internal/http/router/api_group.go
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
)

func ApiGroupRoute(app *fiber.App, h *handler.ApiGroupHandler) {
	apigroup := app.Group("/api-groups")

	apigroup.Use(middleware.Auth())

	apigroup.Post("/", h.CreateNewApiGroupHandler) // Create
	apigroup.Get("/:id", h.GetApiGroupHandler)     // Get One
	apigroup.Put("/", h.UpdateApiGroupHandler)     // Update
	apigroup.Delete("/", h.DeleteApiGroupHandler)  // Delete
}
