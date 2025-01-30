package playlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

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

func (service *PlaylistService) CreatePlaylist(body PlaylistRequest, token string) (CreatePlaylistResponse, error) {
	user, err := service.userService.GetCurrentUserProfile(token)
	if err != nil {
		return CreatePlaylistResponse{}, err
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", url.PathEscape(user.ID))

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

func (service *PlaylistService) GetURIsLikedTracks(body PlaylistRequest, token string, offset int) (URIs []string, total int, err error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?offset=%d&limit=50", offset)

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

func (service *PlaylistService) AddTracksToPlaylist(body PlaylistRequest, token string, URIs []string, playlistId CreatePlaylistResponse, position int) error {
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistId.Id)

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
	var wg sync.WaitGroup

	playlistId, err := service.CreatePlaylist(body, token)
	if err != nil {
		return PlaylistResponse{}, err
	}

	_, totalURIs, err := service.GetURIsLikedTracks(body, token, 0)
	if err != nil {
		return PlaylistResponse{}, err
	}

	for offset := 0; offset < totalURIs; offset += 50 {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()

			URIs, _, err := service.GetURIsLikedTracks(body, token, offset)
			if err != nil {
				return
			}

			err = service.AddTracksToPlaylist(body, token, URIs, playlistId, offset)
			if err != nil {
				return
			}
		}(offset)
	}

	wg.Wait()

	return PlaylistResponse{
		PlaylistId: playlistId.Id,
	}, err
}
