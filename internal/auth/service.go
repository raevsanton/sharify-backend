package auth

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/pkg/codec"
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

	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.Auth.ClientId, config.Auth.ClientSecret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to send request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		tokens, _ := io.ReadAll(res.Body)
		return AuthResponse{}, fmt.Errorf("received non-OK status: %d, body: %s", res.StatusCode, tokens)
	}

	tokens, err := codec.Decode[AuthResponse](res.Body)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("wrong response body: %s", tokens)
	}

	return tokens, nil
}
