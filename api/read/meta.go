package read_api

import (
	"net/http"

	"dev/bluebasooo/video-platform/service"
	"encoding/json"

	"github.com/gorilla/mux"
)

func ReadFileMeta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	meta, err := service.GetFileMeta(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(meta)
}

func VideoPreview(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("preview"))
}

func InitFileMetaApi(router *mux.Router) {
	fileApi := router.PathPrefix("/file").Subrouter()
	fileApi.HandleFunc("/{id}/meta", ReadFileMeta).Methods("GET")
	fileApi.HandleFunc("/{id}/preview", VideoPreview).Methods("GET")
}
