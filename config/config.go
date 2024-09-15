package config

import (
	"log"
	"os"
)

var (
	JwtSecret string
)

func LoadConfig() {
	JwtSecret = os.Getenv("JWT_SECRET")
	if JwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
}
