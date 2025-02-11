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

func App(conf *configs.Config) http.Handler {
	router := http.NewServeMux()

	// Services
	userService := user.NewUseService()
	authService := auth.NewAuthService(userService)
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
	conf := configs.LoadConfig()
	app := App(conf)

	server := http.Server{
		Addr:    ":" + conf.Port,
		Handler: app,
	}
	fmt.Println("Server is listening on port " + conf.Port)
	server.ListenAndServe()
}
