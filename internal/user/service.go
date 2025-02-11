package user

import (
	"net/http"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/pkg/req"
)

type UserService struct {
}

func NewUseService() *UserService {
	return &UserService{}
}

func (service *UserService) GetCurrentUserProfile(token string, config *configs.Config) (CurrentUserResponse, error) {
	r, err := http.NewRequest(http.MethodGet, config.Spotify.ApiUrl+"/me", nil)
	if err != nil {
		return CurrentUserResponse{}, err
	}

	r.Header.Set("Authorization", "Bearer "+token)
	return req.DoRequest[CurrentUserResponse](r, http.StatusOK)
}
