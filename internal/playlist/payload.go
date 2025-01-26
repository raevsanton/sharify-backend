package playlist

type PlaylistRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Public      bool   `json:"is_public"`
}

type PlaylistResponse struct {
	PlaylistId string `json:"playlist_id"`
}

type LikedTracksResponse struct {
	Href   string `json:"href"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Total  int    `json:"total"`
	Items  []Item `json:"items"`
}

type Item struct {
	AddedAt string `json:"added_at"`
	Track   Track  `json:"track"`
}

type Track struct {
	Uri string `json:"uri"`
	// more fields
}

type CreatePlaylistResponse struct {
	Id string `json:"id"`
	// more fields
}

type AddTracksToPlaylistRequest struct {
	URIs     []string `json:"uris" validate:"required"`
	Position int      `json:"position" validate:"required"`
}

type AddTracksToPlaylistResponse struct {
	SnapshotId string `json:"snapshot_id"`
}
