package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

func Register(app *fiber.App, authHandler *handler.AuthHandler, profileHandler *handler.UserProfileHandler,
	apiGroupHandler *handler.ApiGroupHandler, apiHandler *handler.ApiHandler, rbacHandler *handler.RbacHandler,
	rbacService *services.RbacService) {
	app.Get("/health", handler.HealthCheck)
	AuthRoutes(app, authHandler)
	ProfileRoute(app, profileHandler)
	ApiGroupRoute(app, apiGroupHandler)
	ApiRoutes(app, apiHandler, rbacService)
	RbacRoute(app, rbacHandler)
}
