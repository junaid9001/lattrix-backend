package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

func InviteRoutes(
	app *fiber.App, inviteHandler *handler.InvitationHandler,
	notifHandler *handler.NotificationHandler,
	rbacService *services.RbacService) {

	api := app.Group("/api")
	api.Use(middleware.Auth())

	//send invitation
	api.Post("/invitations/send", inviteHandler.SendInvite)
	//accept invitation by token
	api.Post("/invitations/accept", inviteHandler.AcceptInvitation)

	api.Get("/notifications", notifHandler.GetMyNotifications)

	wsNotif := api.Group("/workspace/notifications")
	wsNotif.Get("/", notifHandler.WorkspaceNotifications)

	wsNotif.Delete("/:notificationId",
		middleware.RequirePermission(rbacService, "role:superadmin"),
		notifHandler.DeleteWorkspaceNotification,
	)
}
