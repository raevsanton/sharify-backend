package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth AuthConfig
}

type AuthConfig struct {
	ClientId     string
	ClientSecret string
	ClientUrl    string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Auth: AuthConfig{
			ClientId:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			ClientUrl:    os.Getenv("CLIENT_URL"),
		},
	}
}
