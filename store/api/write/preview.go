package write_api

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateVideoPreview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userId := r.Header.Get("user_id")

	var createVideoPreviewDto dto.CreateVideoPreviewDto
	json.NewDecoder(r.Body).Decode(&createVideoPreviewDto)

	err := service.CreateVideoPreview(id, userId, &createVideoPreviewDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func InitWritePreviewApi(router *mux.Router) {
	fileApi := router.PathPrefix("/preview").Subrouter()
	fileApi.HandleFunc("/{id}", CreateVideoPreview).Methods("POST")
}
