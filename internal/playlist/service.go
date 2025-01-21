package playlist

import (
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

func (service *PlaylistService) CreatePlaylist(body PlaylistRequest) (id CreatePlaylistResponse, err error) {
	user, err := service.userService.GetCurrentUserProfile(user.CurrentUserRequest{})
	if err != nil {
		return CreatePlaylistResponse{}, err
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", url.PathEscape(user.ID))

	req, err := http.NewRequest(http.MethodPost, url, body)
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
		return LikedTracksResponse{}, 0, fmt.Errorf("received non-OK status: %d, body: %s", res.StatusCode, tracks)
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
