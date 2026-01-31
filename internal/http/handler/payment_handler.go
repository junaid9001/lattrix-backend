package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/junaid9001/lattrix-backend/internal/config"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"github.com/junaid9001/lattrix-backend/internal/services"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
	userRepo       repository.UserRepository
	apiRepo        repository.ApiRepository
	cfg            *config.Config
}

func NewPaymentHandler(paymentService *services.PaymentService, userRepo repository.UserRepository, apiRepo repository.ApiRepository, cfg *config.Config) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService, userRepo: userRepo, apiRepo: apiRepo, cfg: cfg}
}

// creates a checkout session and redirect user in frontend
func (h *PaymentHandler) CreateSession(c *fiber.Ctx) error {
	var req struct {
		Plan string `json:"plan"`
	}

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	userIDInt := c.Locals("userID").(int)
	userID := uint(userIDInt)

	url, err := h.paymentService.CreateCheckoutSession(userID, req.Plan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"url": url})

}

func (h *PaymentHandler) HandleWebHook(c *fiber.Ctx) error {
	log.Println(" Webhook endpoint hit!")
	signHeader := c.Get("Stripe-Signature")

	body := c.Body()

	event, err := webhook.ConstructEventWithOptions(body, signHeader, h.cfg.STRIPE_WEBHOOK_SECRET, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		log.Printf("  Webhook signature verification failed: %v\n", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid signature")
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession

		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing JSON")
		}

		h.HandleCheckoutSession(&session)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error parsing JSON")
		}
		h.handleSubscriptionExpired(&subscription)

	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *PaymentHandler) HandleCheckoutSession(sess *stripe.CheckoutSession) {
	userIDStr := sess.Metadata["user_id"]
	planTypeStr := sess.Metadata["plan_type"]

	if userIDStr == "" {
		log.Println(" Missing user_id in session metadata. Note: 'stripe trigger' sends fake data without metadata.")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		log.Printf(" Failed to parse user_id: %v", err)
		return
	}

	var customerID string
	if sess.Customer != nil {
		customerID = sess.Customer.ID
	}

	fmt.Printf(" Payment Success for User %d! Switching to %s\n", userID, planTypeStr)

	updates := map[string]interface{}{
		"plan":                models.PlanType(planTypeStr),
		"stripe_customer_id":  customerID,
		"subscription_status": "active",
	}

	if _, err := h.userRepo.UpdateProfile(uint(userID), updates); err != nil {
		log.Printf(" Failed to update user profile: %v", err)
	}

}

func (h *PaymentHandler) handleSubscriptionExpired(sub *stripe.Subscription) {
	userIDStr := sub.Metadata["user_id"]
	userID64, _ := strconv.ParseUint(userIDStr, 10, 64)
	userID := uint(userID64)

	fmt.Printf(" Downgrading User %d to FREE\n", userID)

	h.userRepo.UpdateProfile(userID, map[string]interface{}{
		"plan":                models.PlanFree,
		"subscription_status": "inactive",
	})

	freeRules := models.PlanRules[models.PlanFree]
	h.apiRepo.EnforcePlanLimits(userID, freeRules.MaxApis, freeRules.MinIntervel)
}
