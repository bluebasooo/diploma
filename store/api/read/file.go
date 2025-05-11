package read_api

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"dev/bluebasooo/video-platform/service"

	"github.com/gorilla/mux"
)

func ReadFilePart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId := r.Header.Get("user_id")
	id := vars["id"]
	part := vars["part"]

	body := bytes.Buffer{}
	answer := multipart.NewWriter(&body)
	chunks, err := service.GetFilePartInterval(userId, id, part)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for hash, chunk := range chunks {
		chunkWriter, err := answer.CreateFormFile(hash, hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		chunkWriter.Write(chunk)
	}

	answer.Close()

	w.Header().Set("Content-Type", answer.FormDataContentType())
	w.Header().Set("Content-Length", fmt.Sprintf("%d", body.Len()))
	w.Write(body.Bytes())
}

func InitFileReadApi(router *mux.Router) {
	fileAPI := router.PathPrefix("/file").Subrouter()
	fileAPI.HandleFunc("/{id}/{part}", ReadFilePart).Methods("GET")
}
