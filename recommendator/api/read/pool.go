package read

import (
	"dev/bluebasooo/video-recomendator/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPoolVideos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	poolId := vars["poolId"]
	page := vars["page"]
	pageSize := vars["pageSize"]

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}

	pool := service.GetPagedVideoPool(poolId, pageInt, pageSizeInt)

	json.NewEncoder(w).Encode(pool)
}

func InitPoolApi(mux *mux.Router) {
	poolRouter := mux.PathPrefix("/pool").Subrouter()
	poolRouter.HandleFunc("/{poolId}", GetPoolVideos).Methods("GET")
}
