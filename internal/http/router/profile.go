package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

func ProfileRoute(app *fiber.App, profileHandler *handler.UserProfileHandler, rabcService *services.RbacService) {
	user := app.Group("/user")
	user.Use(middleware.Auth())
	user.Get("/profile", profileHandler.GetProfile)
	user.Put("/profile", profileHandler.UpdateProfileByID)
	//all users in a workspace
	app.Get("/users", middleware.Auth(), middleware.RequirePermission(rabcService, "role:superadmin"), profileHandler.GetWorkspaceUsers)
}
