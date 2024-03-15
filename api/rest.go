package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(newRouteHandler, redirectHandler func(http.ResponseWriter, *http.Request)) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/new_short_url", handleNewShortUrl).Methods("POST")
	r.PathPrefix("/").HandlerFunc(handleOtherPaths)
	return r
}

func handleNewShortUrl(w http.ResponseWriter, r *http.Request) {
	// Handle new short URL creation here
}

func handleOtherPaths(w http.ResponseWriter, r *http.Request) {
	// Handle all other paths here
}
