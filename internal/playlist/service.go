package playlist

import (
	"fmt"
	"net/http"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/pkg/codec"
)

type PlaylistService struct {
}

func NewPlaylistService() *PlaylistService {
	return &PlaylistService{}
}

func (service *PlaylistService) CreatePlaylist(body PlaylistRequest) (URIs []string, total int, err error) {

}

func (service *PlaylistService) GetURIsLikedTracks(body PlaylistRequest, offset int) (URIs []string, total int, err error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?offset=%d&limit=50", offset)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Authorization", "Bearer "+body.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch liked tracks: %w", err)
	}

	tracks, err := codec.Decode[LikedTracksResponse](res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get liked tracks: %w", err)
	}

	var uris []string
	for _, item := range tracks.Items {
		uris = append(uris, item.Track.Uri)
	}

	return uris, tracks.Total, nil

}

func (service *PlaylistService) GeneratePlaylist(body PlaylistRequest, config *configs.Config) (PlaylistResponse, error) {
	var accumulatedURIs []string
	totalURIs := 0
	offset := 0

	for {
		URIs, total, err := service.GetURIsLikedTracks(body, offset)
		if err != nil {
			return PlaylistResponse{}, fmt.Errorf("failed to get liked tracks URIs: %w", err)
		}
		if totalURIs == 0 {
			totalURIs = total
		}
		offset += 50
		accumulatedURIs = append(accumulatedURIs, URIs...)

		if totalURIs > len(accumulatedURIs) {
			break
		}
	}

	return PlaylistResponse{}, fmt.Errorf("failed to get liked tracks: %w", nil)
}
