package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AppConfig *Config

func Load() {

	godotenv.Load()

	secret := os.Getenv("JWT_KEY")
	if secret == "" {
		log.Fatal("JWT_KEY not set")
	}

	email := os.Getenv("EMAIL")
	emailPass := os.Getenv("EMAIL_PASS")
	stripeSecret := os.Getenv("STRIPE_SECRET_KEY")
	stripePricePro := os.Getenv("STRIPE_PRICE_PRO")
	stripePriceAgency := os.Getenv("STRIPE_PRICE_AGENCY")
	stripeWebHookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	AppConfig = &Config{
		JWT_KEY:               []byte(secret),
		EMAIL:                 email,
		EMAIL_PASSWORD:        emailPass,
		STRIPE_SECRET_KEY:     stripeSecret,
		STRIPE_PRICE_PRO:      stripePricePro,
		STRIPE_PRICE_AGENCY:   stripePriceAgency,
		FRONTEND_URL:          "http://localhost:5173",
		STRIPE_WEBHOOK_SECRET: stripeWebHookSecret,
	}
}

//stripe listen --forward-to localhost:8080/webhooks/stripe
