package application

import (
	db "ci/dbadapters"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type App struct {
	DBAdapter       db.DBAdapter
	Counter         int
	WorkerId        string
	requestsCounter *prometheus.CounterVec
}

func NewApp(dbAdapter db.DBAdapter, workerId string) *App {
	counter, err := dbAdapter.GetCounter()
	if err != nil {
		counter = 0
		log.Printf("Failed to get counter: %v, starting over at count 0", err)
	}
	requestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of requests",
		}, []string{"path"},
	)
	prometheus.MustRegister(requestsCounter)
	return &App{
		DBAdapter:       dbAdapter,
		Counter:         counter,
		WorkerId:        workerId,
		requestsCounter: requestsCounter,
	}
}

func (app *App) AddNewShortUrl(fullUrl string, expiresAt *time.Time) (string, error) {
	tinyUrl := app.generateShortUrl()
	return tinyUrl, app.DBAdapter.InsertNewShortUrl(tinyUrl, fullUrl, expiresAt)
}

func (app *App) GetFullUrl(url string) (string, error) {
	fullUrl, expiryDate, err := app.DBAdapter.GetFullUrl(url)
	if err != nil {
		return "", err
	}
	app.requestsCounter.WithLabelValues(url).Inc()

	if expiryDate != nil && expiryDate.Before(time.Now()) {
		return "", errors.New("url has expired")
	}
	return fullUrl, nil
}

func (app *App) DeleteShortUrl(url string) error {
	return app.DBAdapter.DeleteShortUrl(url)
}

func base62Encode(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(base62Chars[n%62]) + result
		n = n / 62
	}
	return result
}

func (app *App) generateShortUrl() string {
	count := app.Counter
	app.Counter++
	go func() {
		err := app.DBAdapter.UpdateCounter(app.Counter)
		if err != nil {
			log.Printf("Failed to update counter: %v", err)
		}
	}()
	shortUrl := base62Encode(count)
	shortUrl = fmt.Sprintf("%06s", shortUrl)
	return app.WorkerId + shortUrl
}
