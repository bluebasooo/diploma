package read_api

import (
	"context"
	"dev/bluebasooo/video-platform/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoId := vars["videoId"]
	pageNum := vars["pageNum"]
	pageSize := vars["pageSize"]

	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		pageNumInt = 0
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}

	comments, err := service.GetComments(context.Background(), videoId, pageNumInt, pageSizeInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
	w.WriteHeader(http.StatusOK)
}

func InitCommentsReadApi(router *mux.Router) {
	subRouter := router.PathPrefix("/comments").Subrouter()
	subRouter.HandleFunc("/{videoId}/", GetComments).Methods("GET")
}
