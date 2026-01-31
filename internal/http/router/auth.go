package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
)

func AuthRoutes(app *fiber.App, authHandler *handler.AuthHandler) {
	auth := app.Group("/auth")
	auth.Post("/signup", authHandler.Signup)
	auth.Post("/login", authHandler.Login)

	auth.Post("/select-workspace", authHandler.SelectWorkspace)
	auth.Post("/workspace", middleware.Auth(), authHandler.CreateWorkspace)
	auth.Get("/workspaces", middleware.Auth(), authHandler.GetUserWorkspaces)

	auth.Get("/refresh", authHandler.Refresh)
	auth.Get("/me", middleware.Auth(), authHandler.Me)
	auth.Get("/logout", middleware.Auth(), authHandler.Logout)

}
