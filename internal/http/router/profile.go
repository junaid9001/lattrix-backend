package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
)

func ProfileRoute(app *fiber.App, profileHandler *handler.UserProfileHandler) {
	user := app.Group("/user")
	user.Use(middleware.Auth())
	user.Get("/profile", profileHandler.GetProfile)
	user.Put("/profile", profileHandler.UpdateProfileByID)
}
