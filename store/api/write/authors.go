package write_api

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author dto.CreateAuthorDto
	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = service.CreateAuthor(&author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func InitAuthorsApi(router *mux.Router) {
	authorsAPI := router.PathPrefix("/authors").Subrouter()
	authorsAPI.HandleFunc("/create", CreateAuthor).Methods("POST")
}
