package api

import (
	"ci/application"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type UrlShortenRequest struct {
	FullUrl   string `json:"full_url"`
	ExpiresAt string `json:"expires_at"`
}

func SetupRoutes(app *application.App) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/new_short_url", handleNewShortUrl(app)).Methods("POST")
	r.HandleFunc("/{shorturl}", handleDeleteShortUrl(app)).Methods("Delete")
	r.PathPrefix("/").HandlerFunc(handleOtherPaths(app))
	return r
}

func handleNewShortUrl(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		currentShortenRequest := UrlShortenRequest{}

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(bodyBytes, &currentShortenRequest)
		if err != nil {
			http.Error(w, "Failed to unmarshal request body", http.StatusInternalServerError)
			return
		}

		var expiryDate *time.Time
		if currentShortenRequest.ExpiresAt != "" {
			tmpExpiryDate, err := time.Parse(time.RFC3339, currentShortenRequest.ExpiresAt)
			if err != nil {
				http.Error(w, "Invalid expiry date", http.StatusBadRequest)
				return
			}
			expiryDate = &tmpExpiryDate
		}

		tinyUrl, err := app.AddNewShortUrl(currentShortenRequest.FullUrl, expiryDate)
		if err != nil {
			http.Error(w, "Failed to add new short url", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(tinyUrl))
	}
}

func handleDeleteShortUrl(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.DeleteShortUrl(r.URL.Path)
		if err != nil {
			http.Error(w, "Failed to delete short url", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleOtherPaths(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) < 2 {
			http.Error(w, "Invalid url", http.StatusBadRequest)
			return
		}
		fullUrl, err := app.GetFullUrl(r.URL.Path[1:])
		if err != nil {
			http.Error(w, "Failed to get full url", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fullUrl, http.StatusMovedPermanently)
	}
}
