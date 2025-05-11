package write_api

import (
	"context"
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var dto dto.CommentDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = service.CreateComment(context.Background(), &dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func InitWriteCommentsApi(router *mux.Router) {
	subRouter := router.PathPrefix("/comments").Subrouter()
	subRouter.HandleFunc("/", CreateComment).Methods("POST")
}
