package read_api

import (
	"dev/bluebasooo/video-platform/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func FindVideos(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	videos, err := service.FindVideos(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(videos)
}

func InitSearchApi(router *mux.Router) {
	searchApi := router.PathPrefix("/search").Subrouter()
	searchApi.HandleFunc("/videos", FindVideos).Methods("GET")
}
