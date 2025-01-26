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

func (service *PlaylistService) AddTracksToPlaylist(body PlaylistRequest, token string, hundredURIs []string, playlistId CreatePlaylistResponse, position int) error {
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistId.Id)

	bodyRequest := AddTracksToPlaylistRequest{
		URIs:     hundredURIs,
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
	var accumulatedURIs []string
	totalURIs := 0
	offset := 0

	for {
		URIs, total, err := service.GetURIsLikedTracks(body, token, offset)
		if err != nil {
			return PlaylistResponse{}, err
		}
		if totalURIs == 0 {
			totalURIs = total
		}
		offset += 50
		accumulatedURIs = append(accumulatedURIs, URIs...)

		if len(accumulatedURIs) >= totalURIs {
			break
		}
	}

	playlistId, err := service.CreatePlaylist(body, token)
	if err != nil {
		return PlaylistResponse{}, err
	}

	position := 0

	for i := 0; i < len(accumulatedURIs); i += 100 {
		hundredURIs := accumulatedURIs[i:min(i+100, len(accumulatedURIs))]
		err := service.AddTracksToPlaylist(body, token, hundredURIs, playlistId, position)
		if err != nil {
			return PlaylistResponse{}, err
		}
		position += 100
	}

	return PlaylistResponse{
		PlaylistId: playlistId.Id,
	}, err
}
