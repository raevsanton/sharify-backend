package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	Spotify     SpotifyConfig
	Auth        AuthConfig
}

type SpotifyConfig struct {
	ApiUrl  string
	AuthUrl string
}

type AuthConfig struct {
	ClientId     string
	ClientSecret string
	ClientUrl    string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	return &Config{
		Port:        os.Getenv("PORT"),
		Environment: environment,
		Spotify: SpotifyConfig{
			ApiUrl:  os.Getenv("SPOTIFY_API_URL"),
			AuthUrl: os.Getenv("SPOTIFY_AUTH_URL"),
		},
		Auth: AuthConfig{
			ClientId:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			ClientUrl:    os.Getenv("CLIENT_URL"),
		},
	}
}
