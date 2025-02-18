package playlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/internal/user"
	"github.com/raevsanton/sharify-backend/pkg/req"
)

type PlaylistService struct {
	userService *user.UserService
}

func NewPlaylistService(userService *user.UserService) *PlaylistService {
	return &PlaylistService{
		userService: userService,
	}
}

func (service *PlaylistService) CreatePlaylist(body PlaylistRequest, config *configs.Config, token string) (CreatePlaylistResponse, error) {
	user, err := service.userService.GetCurrentUserProfile(token, config)
	if err != nil {
		return CreatePlaylistResponse{}, err
	}

	url := fmt.Sprintf("%s/users/%s/playlists", config.Spotify.ApiUrl, url.PathEscape(user.ID))

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return CreatePlaylistResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return CreatePlaylistResponse{}, err
	}

	r.Header.Set("Authorization", "Bearer "+token)
	return req.DoRequest[CreatePlaylistResponse](r, http.StatusCreated)
}

func (service *PlaylistService) GetURIsLikedTracks(body PlaylistRequest, config *configs.Config, token string, offset int) (URIs []string, total int, err error) {
	url := fmt.Sprintf("%s/me/tracks?offset=%d&limit=50", config.Spotify.ApiUrl, offset)

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	r.Header.Set("Authorization", "Bearer "+token)

	tracks, err := req.DoRequest[LikedTracksResponse](r, http.StatusOK)
	if err != nil {
		return nil, 0, err
	}

	var uris []string
	for _, item := range tracks.Items {
		uris = append(uris, item.Track.Uri)
	}

	return uris, tracks.Total, nil
}

func (service *PlaylistService) AddTracksToPlaylist(body PlaylistRequest, config *configs.Config, token string, URIs []string, playlistId CreatePlaylistResponse, position int) error {
	url := fmt.Sprintf("%s/playlists/%s/tracks", config.Spotify.ApiUrl, playlistId.Id)

	bodyRequest := AddTracksToPlaylistRequest{
		URIs:     URIs,
		Position: position,
	}

	bodyJson, err := json.Marshal(bodyRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyJson))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch adding tracks to playlist: %w", err)
	}

	return nil
}

func (service *PlaylistService) GeneratePlaylist(body PlaylistRequest, config *configs.Config, token string) (PlaylistResponse, error) {
	playlistId, err := service.CreatePlaylist(body, config, token)
	if err != nil {
		return PlaylistResponse{}, err
	}

	_, totalURIs, err := service.GetURIsLikedTracks(body, config, token, 0)
	if err != nil {
		return PlaylistResponse{}, err
	}

	for offset := 0; offset < totalURIs; offset += 50 {
		URIs, _, err := service.GetURIsLikedTracks(body, config, token, offset)
		if err != nil {
			return PlaylistResponse{}, err
		}

		err = service.AddTracksToPlaylist(body, config, token, URIs, playlistId, offset)
		if err != nil {
			return PlaylistResponse{}, err
		}
	}

	return PlaylistResponse{
		PlaylistId: playlistId.Id,
	}, err
}
