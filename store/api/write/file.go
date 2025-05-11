package write_api

import (
	"dev/bluebasooo/video-platform/api/dto"
	"dev/bluebasooo/video-platform/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func WriteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["taskId"]
	userId := r.Header.Get("user_id")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	resultChanel := make(chan *channelResult, len(r.MultipartForm.File))

	for key, value := range r.MultipartForm.File {
		for _, fileHeader := range value {

			file, _ := fileHeader.Open()
			defer file.Close()
			bytes, _ := io.ReadAll(file)

			wg.Add(1)
			go func(id string, wg *sync.WaitGroup, resultChan chan<- *channelResult) {
				defer wg.Done()
				err := service.Write(taskID, id, bytes, userId)
				resultChan <- &channelResult{
					id:  id,
					err: err,
				}

			}(key, &wg, resultChanel)
		}
	}

	wg.Wait()
	close(resultChanel)

	failed := make([]dto.WriteResultDto, 0)
	for result := range resultChanel {
		if result.err != nil {
			failed = append(failed, dto.WriteResultDto{
				ID:    result.id,
				Error: result.err.Error(),
			})
		}
	}

	json.NewEncoder(w).Encode(&dto.WritePartsResultDto{
		Results: failed,
	})

	w.WriteHeader(http.StatusOK)
}

type channelResult struct {
	id  string
	err error
}

func GeneratePlan(w http.ResponseWriter, r *http.Request) {
	var planDto dto.FileMetaPlanDto
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

func InitWriteFileApi(router *mux.Router) {
	fileAPI := router.PathPrefix("/file").Subrouter()
	fileAPI.HandleFunc("/write/plan", GeneratePlan).Methods("POST")
	fileAPI.HandleFunc("/write/{taskId}", WriteFile).Methods("POST")
}
