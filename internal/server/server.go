package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/junaid9001/lattrix-backend/internal/config"
	"github.com/junaid9001/lattrix-backend/internal/consumer"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/router"
	"github.com/junaid9001/lattrix-backend/internal/infra"
	"github.com/junaid9001/lattrix-backend/internal/infra/repo"
	"github.com/junaid9001/lattrix-backend/internal/publisher"
	"github.com/junaid9001/lattrix-backend/internal/services"
	"github.com/junaid9001/lattrix-backend/internal/worker"
)

func Start() {
	config.Load()
	db := infra.ConnectDB()

	//message broker
	kafkaWriter := publisher.NewKafkaWriter(publisher.Kafkaconfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "api-jobs",
	})
	defer kafkaWriter.Close()

	kafkaReader := consumer.NewKafkaConsumer(consumer.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "api-jobs",
		GroupID: "lattrix-workers",
	})
	defer kafkaReader.Close()

	kafkaPublisher := publisher.NewkafkaPublisher(kafkaWriter)
	KafkaConsumer := consumer.NewKafkaJobConsumer(kafkaReader)

	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("shutdown signal received")
		cancel()
	}()

	//repos
	userRepo := repo.NewUserRepository(db)
	apiGroupRepo := repo.NewApiGroupRepository(db)
	apiRepo := repo.NewApiRepo(db)
	rbacRepo := repo.NewRbacRepo(db)
	inviteRepo := repo.NewInvitationRepo(db)
	notificationRepo := repo.NewNotificationRepo(db)
	workNotiRepo := repo.NewWorkspaceNotificationRepository(db)

	workNotiService := services.NewWorkspaceNotiService(workNotiRepo)

	schedular := worker.NewSchedular(apiRepo, kafkaPublisher)
	go schedular.Start(ctx)

	Workerconsumer := worker.NewConsumer(KafkaConsumer, workNotiService, db)
	const workerCount = 24
	for i := 0; i < workerCount; i++ {
		go Workerconsumer.Start(ctx)
	}

	//services
	authService := services.NewAuthSevice(userRepo, apiGroupRepo, rbacRepo, db)
	profileService := services.NewProfileService(userRepo)
	apiGroupService := services.NewApiGroupService(apiGroupRepo, userRepo)
	apiService := services.NewApiService(apiRepo, userRepo, db)
	rbacService := services.NewRbacService(rbacRepo, db)
	inviteService := services.NewInvitationService(db, inviteRepo, notificationRepo, rbacRepo, userRepo)
	paymentService := services.NewPaymentService(userRepo, config.AppConfig)
	adminService := services.NewAdminService(db, kafkaWriter, config.AppConfig)

	//handlers
	paymentHandler := handler.NewPaymentHandler(paymentService, userRepo, apiRepo, config.AppConfig)
	authHandler := handler.NewAuthHandler(authService, rbacService, profileService)
	profileHandler := handler.NewProfileHandler(profileService)
	apiGroupHandler := handler.NewApiGroupHandler(apiGroupService)
	apiHandler := handler.NewApiHandler(apiService)
	rbacHandler := handler.NewRbacHandler(rbacService)
	adminHandler := handler.NewAdminHandler(adminService)

	inviteHandler := handler.NewInvitationHandler(inviteService)
	notificationHandler := handler.NewNotificationHandler(notificationRepo, workNotiRepo)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	router.Register(app, authHandler, profileHandler, apiGroupHandler, apiHandler, rbacHandler, rbacService)
	router.InviteRoutes(app, inviteHandler, notificationHandler, rbacService)
	router.PaymentRoutes(app, paymentHandler)
	router.AdminRoutes(app, adminHandler, db)

	log.Fatal(app.Listen(":8080"))
}
