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
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.Handle("/robots.txt", http.FileServer(http.Dir(ASSETS_DIR)))
	mux.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(ASSETS_DIR))))

	mux.Handle("/logs/{id:[0-9]+}", ServerHandlerFunc{server: server, f: showLogPage})
	mux.Handle("/", ServerHandlerFunc{server: server, f: indexPage})
	return mux
}
