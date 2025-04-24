package read_api

import (
	"dev/bluebasooo/video-platform/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func VideoPreview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	preview, err := service.GetVideoPreview(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(preview)
}

func GetMainPageVideoPreviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"] // TODO: get from auth

	preview, err := service.GetMainPageVideoPreviews(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(preview)
}

func InitPreviewApi(router *mux.Router) {
	fileApi := router.PathPrefix("/preview").Subrouter()
	fileApi.HandleFunc("/{id}/preview", VideoPreview).Methods("GET")
	fileApi.HandleFunc("/main-page/{userID}", GetMainPageVideoPreviews).Methods("GET")
}
