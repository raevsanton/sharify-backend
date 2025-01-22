package playlist

import (
	"net/http"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/pkg/req"
	"github.com/raevsanton/sharify-backend/pkg/res"
)

type PlaylistHandlerDeps struct {
	*configs.Config
	*PlaylistService
}

type PlaylistHandler struct {
	*configs.Config
	*PlaylistService
}

func NewPlaylistHandler(router *http.ServeMux, deps PlaylistHandlerDeps) {
	handler := &PlaylistHandler{
		Config:          deps.Config,
		PlaylistService: deps.PlaylistService,
	}
	router.HandleFunc("POST /playlist", handler.Playlist(deps.Config))
}

func (handler *PlaylistHandler) Playlist(config *configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[PlaylistRequest](&w, r)
		if err != nil {
			return
		}

		playlistId, err := handler.PlaylistService.GeneratePlaylist(*body, config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, playlistId, http.StatusOK)
	}
}
