package playlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/internal/user"
	"github.com/raevsanton/sharify-backend/pkg/codec"
)

type PlaylistService struct {
	userService *user.UserService
}

func NewPlaylistService(userService *user.UserService) *PlaylistService {
	return &PlaylistService{
		userService: userService,
	}
}

func (service *PlaylistService) CreatePlaylist(body PlaylistRequest) (CreatePlaylistResponse, error) {
	user, err := service.userService.GetCurrentUserProfile(user.CurrentUserRequest{})
	if err != nil {
		return CreatePlaylistResponse{}, err
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", url.PathEscape(user.ID))

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return CreatePlaylistResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return CreatePlaylistResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+body.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return CreatePlaylistResponse{}, fmt.Errorf("failed to send request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		tokens, _ := io.ReadAll(res.Body)
		return CreatePlaylistResponse{}, fmt.Errorf("received non-OK status: %d, body: %s", res.StatusCode, tokens)
	}

	playlist, err := codec.Decode[CreatePlaylistResponse](res.Body)
	if err != nil {
		return CreatePlaylistResponse{}, fmt.Errorf("failed to get liked tracks: %w", err)
	}

	return playlist, nil

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

	if res.StatusCode != http.StatusOK {
		tracks, _ := io.ReadAll(res.Body)
		return nil, 0, fmt.Errorf("received non-OK status: %d, body: %s", res.StatusCode, tracks)
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

func (service *PlaylistService) AddTracksToPlaylist(body PlaylistRequest, hundredURIs []string, playlistId CreatePlaylistResponse, position int) error {
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%d/tracks", playlistId.Id)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+body.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch adding tracks to playlist: %w", err)
	}

	return fmt.Errorf("received non-OK status: %d", res.StatusCode)
}

func (service *PlaylistService) GeneratePlaylist(body PlaylistRequest, config *configs.Config) (PlaylistResponse, error) {
	var accumulatedURIs []string
	totalURIs := 0
	offset := 0

	for {
		URIs, total, err := service.GetURIsLikedTracks(body, offset)
		if err != nil {
			return PlaylistResponse{}, err
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

	playlistId, err := service.CreatePlaylist(body)
	if err != nil {
		return PlaylistResponse{}, err
	}

	position := 0

	for i := 0; i < len(accumulatedURIs); i += 100 {
		hundredURIs := accumulatedURIs[i:min(i+100, len(accumulatedURIs))]
		err := service.AddTracksToPlaylist(body, hundredURIs, playlistId, position)
		if err != nil {
			return PlaylistResponse{}, err
		}
		position += 100
	}

	return PlaylistResponse{
		PlaylistId: playlistId.Id,
	}, err
}
