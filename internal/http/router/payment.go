package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/http/handler"
	"github.com/junaid9001/lattrix-backend/internal/http/middleware"
)

func PaymentRoutes(app *fiber.App, h *handler.PaymentHandler) {
	payment := app.Group("/subscription")

	// Create Checkout Link
	payment.Post("/checkout", middleware.Auth(), h.CreateSession)

	// Public: Stripe Webhook
	app.Post("/webhooks/stripe", h.HandleWebHook)
}
