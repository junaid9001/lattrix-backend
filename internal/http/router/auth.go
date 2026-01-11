package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
)

func AuthRoutes(app *fiber.App, authHandler *handler.AuthHandler) {
	auth := app.Group("/auth")
	auth.Post("/signup", authHandler.Signup)
	auth.Post("/login", authHandler.Login)
	auth.Get("/refresh", authHandler.Refresh)

}
