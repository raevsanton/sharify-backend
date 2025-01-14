package auth

import (
	"net/http"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/pkg/req"
	"github.com/raevsanton/sharify-backend/pkg/res"
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
	router.HandleFunc("POST /auth", handler.Auth(deps.Config))
}

func (handler *AuthHandler) Auth(config *configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[AuthRequest](&w, r)
		if err != nil {
			return
		}

		tokens, err := handler.AuthService.GetTokens(body.AuthorizationCode, config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, tokens, http.StatusOK)
	}
}
