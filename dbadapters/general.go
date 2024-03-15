package dpadapters

import (
	"log"
	"strings"
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
	CreateTablesAndStatements(string) error
	InsertNewShortUrl(string, string, *time.Time) error
	GetFullUrl(string) (string, *time.Time, error)
	DeleteShortUrl(string) error
	GetCounter(string) (int, error)
	UpdateCounter(string, int) error
	Close()
	Cleanup()
	printAllTableContents()
}

func ConvertDBType(dbType string) DBType {
	switch strings.ToLower(dbType) {
	case "postgres":
		return Postgres
	case "sqlite":
		return SQLite
	}
	log.Printf("DBType is not supported")
	return -1
}

func NewDBAdapter(dbType DBType, workerId string) DBAdapter {
	switch dbType {
	case Postgres:
		return newPostgresConnection()
	case SQLite:
		return newSQLiteConnection()
	}

	log.Printf("DBType is not supported")
	return nil
}
