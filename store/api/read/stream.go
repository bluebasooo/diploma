package read_api

// func GetRealTimeStream(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	streamID := vars["streamID"]

// 	bytes, err := service.GetRealTimeStream(streamID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write(bytes)
// 	w.WriteHeader(http.StatusOK)
// }

// func GetLastTimeStream(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	streamID := vars["streamID"]

// 	bytes, err := service.GetLastTimeStream(streamID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// func InitStreamApi(router *mux.Router) {
// 	subRouter := router.PathPrefix("/stream").Subrouter()

// 	subRouter.HandleFunc("/{streamID}/real-time", GetRealTimeStream).Methods("GET")
// 	subRouter.HandleFunc("/{streamID}/last-time", GetLastTimeStream).Methods("GET")
// }
