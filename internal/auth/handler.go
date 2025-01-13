package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/pkg/req"
	"github.com/raevsanton/sharify-backend/pkg/res"
)

type AuthHandlerDeps struct {
	*configs.Config
}

type AuthHandler struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth", handler.Auth())
}

func (handler *AuthHandler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[AuthRequest](&w, r)

		if err != nil {
			return
		}

		clientID := os.Getenv("SPOTIFY_CLIENT_ID")
		clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
		redirectURI := os.Getenv("CLIENT_URL")

		data := url.Values{}
		data.Set("grant_type", "authorization_code")
		data.Set("redirect_uri", redirectURI)
		data.Set("code", body.AuthorizationCode)

		req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", nil)
		if err != nil {
			http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(clientID, clientSecret)

		fmt.Println(body)

		data := AuthResponse{
			RefreshToken: "123",
			AccessToken:  "321",
		}

		res.Json(w, data, http.StatusOK)
	}
}
