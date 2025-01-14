package main

import (
	"fmt"
	"net/http"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/internal/auth"
	"github.com/raevsanton/sharify-backend/pkg/middleware"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()

	// Services
	authService := auth.NewAuthService()

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: middleware.CORS(router),
	}
	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
