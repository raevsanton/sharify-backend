package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (service *AuthService) GetTokens(authCode string) (AuthResponse, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := os.Getenv("CLIENT_URL")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)
	data.Set("code", authCode)

	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return AuthResponse{}, fmt.Errorf("received non-OK status: %d, body: %s", resp.StatusCode, body)
	}

	var tokens AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return AuthResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return tokens, nil
}
