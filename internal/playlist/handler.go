package playlist

import (
	"net/http"

	"github.com/raevsanton/sharify-backend/configs"
	"github.com/raevsanton/sharify-backend/pkg/cookie"
	"github.com/raevsanton/sharify-backend/pkg/middleware"
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
	router.Handle("POST /playlist", middleware.IsAuthed(handler.Playlist(deps.Config), deps.Config))
}

func (handler *PlaylistHandler) Playlist(config *configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[PlaylistRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		accessToken, err := cookie.GetCookie(r, "access_token")
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		playlistId, err := handler.PlaylistService.GeneratePlaylist(*body, config, accessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, playlistId, http.StatusOK)
	}
}
