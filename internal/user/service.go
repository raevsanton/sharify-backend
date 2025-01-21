package user

import (
	"fmt"
	"net/http"

	"github.com/raevsanton/sharify-backend/pkg/codec"
)

type UserService struct {
}

func NewUseService() *UserService {
	return &UserService{}
}

func (service *UserService) GetCurrentUserProfile(body CurrentUserRequest) (CurrentUserResponse, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return CurrentUserResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+body.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return CurrentUserResponse{}, fmt.Errorf("failed to fetch user data: %w", err)
	}

	user, err := codec.Decode[CurrentUserResponse](res.Body)
	if err != nil {
		return CurrentUserResponse{}, fmt.Errorf("failed to get user data: %w", err)
	}

	return user, nil
}
