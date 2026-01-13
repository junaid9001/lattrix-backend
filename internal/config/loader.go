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

	AppConfig = &Config{
		JWT_KEY: []byte(secret),
	}
}
