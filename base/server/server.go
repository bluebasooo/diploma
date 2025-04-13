package base_server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type FustestServerUSee struct {
	server *http.Server
	router *mux.Router
}

type RouteInitializer = func(router *mux.Router)

func (fsus *FustestServerUSee) Init(apiInit []RouteInitializer) {
	router := mux.NewRouter()
	fsus.server = &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	fsus.router = router

	for _, initializer := range apiInit {
		fsus.AddRoute(initializer)
	}
}

func (fsus *FustestServerUSee) Start() {
	_ = fsus.server.ListenAndServe()
}

func (fsus *FustestServerUSee) AddRoute(initializer RouteInitializer) {
	initializer(fsus.router)
}
