package playlist

import (
	"github.com/raevsanton/sharify-backend/configs"
)

type PlaylistService struct {
}

func NewPlaylistService() *PlaylistService {
	return &PlaylistService{}
}

func (service *PlaylistService) GeneratePlaylist(body PlaylistRequest, config *configs.Config) (PlaylistResponse, error) {

}
