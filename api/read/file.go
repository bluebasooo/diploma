package read_api

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strings"

	"dev/bluebasooo/video-platform/service"

	"github.com/gorilla/mux"
)

func ReadFilePart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	part := vars["parts"]

	parts := strings.Split(part, ",")

	body := bytes.Buffer{}
	answer := multipart.NewWriter(&body)
	for _, part := range parts {
		chunk, err := service.GetFilePart(id, part)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		chunkWriter, err := answer.CreateFormFile(part, part)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		chunkWriter.Write(chunk)
	}

	answer.Close()

	w.Header().Set("Content-Type", answer.FormDataContentType())
	w.Write(body.Bytes())
}

func InitFileReadApi(router *mux.Router) {
	fileAPI := router.PathPrefix("/file").Subrouter()
	fileAPI.HandleFunc("/{id}/part", ReadFilePart).Methods("GET")
}
