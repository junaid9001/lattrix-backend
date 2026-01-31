package config

type Config struct {
	JWT_KEY        []byte
	EMAIL          string
	EMAIL_PASSWORD string

	STRIPE_SECRET_KEY   string
	STRIPE_PRICE_PRO    string
	STRIPE_PRICE_AGENCY string
	FRONTEND_URL        string

	STRIPE_WEBHOOK_SECRET string
}
