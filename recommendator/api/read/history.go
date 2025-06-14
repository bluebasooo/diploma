package read

import (
	"dev/bluebasooo/video-recomendator/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUserHistoryVideos(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	history, err := service.GetUserHistory(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(history)
}

func InitHistoryApi(mux *mux.Router) {
	poolRouter := mux.PathPrefix("/history").Subrouter()
	poolRouter.HandleFunc("/", GetUserHistoryVideos).Methods("GET")
}
