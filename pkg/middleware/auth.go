package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/internal/auth"
	"github.com/raevsanton/sharify-backend/pkg/cookie"
	"github.com/raevsanton/sharify-backend/pkg/req"
)

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, errAccessToken := cookie.GetCookie(r, "access_token")
		refreshToken, errRefreshToken := cookie.GetCookie(r, "refresh_token")

		if errRefreshToken != nil {
			http.Error(w, "There's no refresh token", http.StatusUnauthorized)
			return
		}

		if errAccessToken != nil {
			tokens, err := getNewTokens(config, refreshToken)

			if err != nil {
				http.Error(w, "Access token haven't refreshed", http.StatusUnauthorized)
				return
			}

			cookie.SetCookie(w, "access_token", tokens.AccessToken, 3600)

			r.AddCookie(&http.Cookie{
				Name:  "access_token",
				Value: tokens.AccessToken,
			})
		}
		next.ServeHTTP(w, r)
	})
}

func getNewTokens(config *configs.Config, token string) (auth.AuthResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", token)

	r, err := http.NewRequest(http.MethodPost, config.Spotify.AuthUrl+"/token", strings.NewReader(data.Encode()))
	if err != nil {
		return auth.AuthResponse{}, err
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetBasicAuth(config.Auth.ClientId, config.Auth.ClientSecret)

	return req.DoRequest[auth.AuthResponse](r, http.StatusOK)
}
