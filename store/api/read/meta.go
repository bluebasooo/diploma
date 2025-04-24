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

func InitFileMetaApi(router *mux.Router) {
	fileApi := router.PathPrefix("/meta").Subrouter()
	fileApi.HandleFunc("/{id}/meta", ReadFileMeta).Methods("GET")
}
