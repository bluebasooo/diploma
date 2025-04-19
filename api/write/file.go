package write_api

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/service"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["taskId"]
	hash := vars["hash"]

	bytes, err := io.ReadAll(r.Body)
	if err != nil { // TODO: handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return // no handle bytes
	}

	err = service.Write(taskID, hash, bytes)
	if err != nil { // TODO: handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return // no handle bytes
	}

	w.WriteHeader(http.StatusOK)
}

func GeneratePlan(w http.ResponseWriter, r *http.Request) {
	var planDto dto.PlainFileMetaDto
	err := json.NewDecoder(r.Body).Decode(&planDto)
	if err != nil { // TODO: handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return // no handle bytes
	}

	plan, err := service.GeneratePlan(&planDto)
	if err != nil { // TODO: handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return // no handle bytes
	}

	json.NewEncoder(w).Encode(plan)
}

func InitFileWriteApi(router *mux.Router) {
	fileAPI := router.PathPrefix("/file").Subrouter()
	fileAPI.HandleFunc("/write/{taskId}/{hash}", WriteFile).Methods("POST")
	fileAPI.HandleFunc("/write/plan", GeneratePlan).Methods("POST")
}
