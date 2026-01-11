package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
)

func InviteRoutes(app *fiber.App, inviteHandler *handler.InvitationHandler, notifHandler *handler.NotificationHandler) {
	api := app.Group("/api")
	api.Use(middleware.Auth())

	api.Post("/invitations/send", inviteHandler.SendInvite)
	api.Post("/invitations/accept", inviteHandler.AcceptInvitation)

	api.Get("/notifications", notifHandler.GetMyNotifications)
}
