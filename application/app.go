package application

import (
	db "ci/dbadapters"
	"errors"
	"time"
)

type App struct {
	DBAdapter db.DBAdapter
	WorkerId  string
}

func NewApp(dbAdapter db.DBAdapter) *App {
	return &App{
		DBAdapter: dbAdapter,
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
	if expiryDate != nil && expiryDate.Before(time.Now()) {
		return "", errors.New("url has expired")
	}

	return fullUrl, nil
}

func (app *App) DeleteShortUrl(url string) error {
	return app.DBAdapter.DeleteShortUrl(url)
}

func (app *App) generateShortUrl() string {
	return app.WorkerId + "shorturl"
}
