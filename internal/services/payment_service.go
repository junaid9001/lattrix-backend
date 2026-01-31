package services

import (
	"fmt"

	"github.com/junaid9001/lattrix-backend/internal/config"
	"github.com/junaid9001/lattrix-backend/internal/domain/repository"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

type PaymentService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewPaymentService(userRepo repository.UserRepository, config *config.Config) *PaymentService {
	stripe.Key = config.STRIPE_SECRET_KEY
	return &PaymentService{userRepo: userRepo, cfg: config}

}

func (s *PaymentService) CreateCheckoutSession(userID uint, planType string) (string, error) {
	user, err := s.userRepo.FindByID(int(userID))
	if err != nil {
		return "", err
	}

	if string(user.Plan) == planType {
		return "", fmt.Errorf("you are already subscribed to the %s plan", planType)
	}

	var priceID string

	switch planType {
	case "PRO":
		priceID = s.cfg.STRIPE_PRICE_PRO
	case "AGENCY":
		priceID = s.cfg.STRIPE_PRICE_AGENCY
	default:
		return "", fmt.Errorf("invalid plan type: %s", planType)
	}

	params := &stripe.CheckoutSessionParams{
		CustomerEmail: &user.Email,

		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{Price: &priceID, Quantity: stripe.Int64(1)},
		},

		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(s.cfg.FRONTEND_URL + "/dashboard?payment=success"),
		CancelURL:  stripe.String(s.cfg.FRONTEND_URL + "/pricing?payment=cancelled"),

		Metadata: map[string]string{
			"user_id":   fmt.Sprintf("%d", userID),
			"plan_type": string(planType),
		},

		ClientReferenceID: stripe.String(fmt.Sprintf("%d", userID)),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id":   fmt.Sprintf("%d", userID),
				"plan_type": string(planType),
			},
		},
	}
	sess, err := session.New(params)
	if err != nil {
		return "", err
	}

	return sess.URL, nil

}
