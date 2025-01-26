package auth

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/pkg/req"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (service *AuthService) GetTokens(authCode string, config *configs.Config) (AuthResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", config.Auth.ClientUrl)
	data.Set("code", authCode)

	r, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return AuthResponse{}, err
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetBasicAuth(config.Auth.ClientId, config.Auth.ClientSecret)

	return req.DoRequest[AuthResponse](r, http.StatusOK)
}
