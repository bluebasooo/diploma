package write

import (
	"dev/bluebasooo/video-recomendator/api/dto"
	"dev/bluebasooo/video-recomendator/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteMetrics(w http.ResponseWriter, r *http.Request) {
	var metrics dto.MetricsDto
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.WriteMetrics(r.Context(), metrics.Metrics)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func InitMetricsApi(mux *mux.Router) {
	metricsRouter := mux.PathPrefix("/metrics").Subrouter()
	metricsRouter.HandleFunc("/", WriteMetrics).Methods("POST")
}
