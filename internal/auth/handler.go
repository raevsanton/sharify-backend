package auth

import (
	"net/http"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/pkg/cookie"
	"github.com/raevsanton/sharify-backend/pkg/req"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth", handler.Authentication(deps.Config))
	router.HandleFunc("GET /auth", handler.TokenValidation(deps.Config))
}

func (handler *AuthHandler) Authentication(config *configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[AuthRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokens, err := handler.AuthService.GetTokens(body.AuthorizationCode, config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cookie.SetCookie(w, "access_token", tokens.AccessToken, 3600)
		cookie.SetCookie(w, "refresh_token", tokens.RefreshToken, 604800)

		w.WriteHeader(http.StatusOK)
	}
}

func (handler *AuthHandler) TokenValidation(config *configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := cookie.GetCookie(r, "access_token")
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		_, err = handler.userService.GetCurrentUserProfile(token, config)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
