package read_api

import (
	"dev/bluebasooo/video-platform/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	author, err := service.GetAuthor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(author)
}
