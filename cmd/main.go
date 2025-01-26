package main

import (
	"fmt"
	"net/http"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/internal/auth"
	"github.com/raevsanton/sharify-backend/internal/playlist"
	"github.com/raevsanton/sharify-backend/internal/user"
	"github.com/raevsanton/sharify-backend/pkg/middleware"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	router := http.NewServeMux()

	// Services
	authService := auth.NewAuthService()
	userService := user.NewUseService()
	playlistService := playlist.NewPlaylistService(userService)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	playlist.NewPlaylistHandler(router, playlist.PlaylistHandlerDeps{
		Config:          conf,
		PlaylistService: playlistService,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
	)

	return stack(router)
}

func main() {
	app := App()
	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}
	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
