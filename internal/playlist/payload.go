package playlist

type PlaylistRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Public       bool   `json:"is_public"`
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type PlaylistResponse struct {
	PlaylistId  string `json:"playlist_id"`
	AccessToken string `json:"access_token"`
}
