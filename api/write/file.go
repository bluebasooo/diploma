package write_api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func WriteFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("write"))
}

func InitFileWriteApi(router *mux.Router) {
	fileApi := router.PathPrefix("/file").Subrouter()
	fileApi.HandleFunc("/write", WriteFile).Methods("POST")
}
