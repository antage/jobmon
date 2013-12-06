package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type ServerHandlerFunc struct {
	server *Server
	f      func(*Server, http.ResponseWriter, *http.Request)
}

func (h ServerHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.f(h.server, w, r)
}

func routes(server *Server) *mux.Router {
	mux := mux.NewRouter()

	// assets
	mux.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(ASSETS_DIR))))

	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.Handle("/robots.txt", http.FileServer(http.Dir(ASSETS_DIR)))
	mux.Handle("/jobs.html", http.FileServer(http.Dir(ASSETS_DIR)))

	// app
	mux.HandleFunc("/", appPage)
	mux.HandleFunc("/logs/{id:[0-9]+}", appPage)

	// json
	mux.Handle("/logs/{id:[0-9]+}.json", ServerHandlerFunc{server: server, f: logEntryJson})
	mux.Handle("/jobs.json", ServerHandlerFunc{server: server, f: jobsListJson})

	return mux
}
