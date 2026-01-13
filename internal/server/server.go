package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/junaid9001/lattrix-backend/internal/config"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/router"
	"github.com/junaid9001/lattrix-backend/internal/infra"
	"github.com/junaid9001/lattrix-backend/internal/infra/repo"
	"github.com/junaid9001/lattrix-backend/internal/services"
)

func Start() {
	config.Load()
	db := infra.ConnectDB()

	userRepo := repo.NewUserRepository(db)
	apiGroupRepo := repo.NewApiGroupRepository(db)
	apiRepo := repo.NewApiRepo(db)
	rbacRepo := repo.NewRbacRepo(db)
	inviteRepo := repo.NewInvitationRepo(db)
	notificationRepo := repo.NewNotificationRepo(db)

	authService := services.NewAuthSevice(userRepo, apiGroupRepo, rbacRepo, db)
	profileService := services.NewProfileService(userRepo)
	apiGroupService := services.NewApiGroupService(apiGroupRepo, userRepo)
	apiService := services.NewApiService(apiRepo)
	rbacService := services.NewRbacService(rbacRepo, db)
	inviteService := services.NewInvitationService(db, inviteRepo, notificationRepo, rbacRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService, rbacService)
	profileHandler := handler.NewProfileHandler(profileService)
	apiGroupHandler := handler.NewApiGroupHandler(apiGroupService)
	apiHandler := handler.NewApiHandler(apiService)
	rbacHandler := handler.NewRbacHandler(rbacService)

	inviteHandler := handler.NewInvitationHandler(inviteService)
	notificationHandler := handler.NewNotificationHandler(notificationRepo)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	router.Register(app, authHandler, profileHandler, apiGroupHandler, apiHandler, rbacHandler, rbacService)
	router.InviteRoutes(app, inviteHandler, notificationHandler)

	log.Fatal(app.Listen(":8080"))
}
