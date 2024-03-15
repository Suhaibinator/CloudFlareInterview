package dpadapters

import (
	"log"
	"time"
)

type DBType int

const (
	Postgres DBType = iota
	SQLite
)

func (dbType DBType) String() string {
	switch dbType {
	case Postgres:
		return "Postgres"
	case SQLite:
		return "SQLite"
	}
	return "Unknown"
}

type DBConfig struct {
	DBType   DBType
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	FilePath string
}

type DBAdapter interface {
	Connect(dbConfig DBConfig) error
	CreateTablesAndStatements() error
	InsertNewShortUrl(url, fullUrl string, expiresAt *time.Time) error
	GetFullUrl(url string) (string, *time.Time, error)
	DeleteShortUrl(url string) error
	GetCounter() (int, error)
	UpdateCounter(int) error
	Close()
	Cleanup()
	printAllTableContents()
}

func NewDBAdapter(dbType DBType) DBAdapter {
	switch dbType {
	case Postgres:
		return newPostgresConnection()
	case SQLite:
		return newSQLiteConnection()
	}

	log.Printf("DBType is not supported")
	return nil
}
