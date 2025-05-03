package write_api

// func CreateStream(w http.ResponseWriter, r *http.Request) {
// 	var streamDto dto.StreamDto
// 	err := json.NewDecoder(r.Body).Decode(&streamDto)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	service.PrepareStream(streamDto.ID)

// 	w.WriteHeader(http.StatusOK)
// }

// func EndStream(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	streamID := vars["streamID"]

// 	service.EndStream(streamID)

// 	w.WriteHeader(http.StatusOK)
// }

// func InitStreamApi(router *mux.Router) {
// 	subRouter := router.PathPrefix("/stream").Subrouter()

// 	subRouter.HandleFunc("/{streamID}/create", CreateStream).Methods("POST")
// 	subRouter.HandleFunc("/{streamID}/end", EndStream).Methods("POST")
// }
