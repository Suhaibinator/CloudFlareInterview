package dpadapters

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteConnection struct {
	mux               *sync.RWMutex
	conn              *sql.DB
	insertNewShortUrl *sql.Stmt
}

func newSQLiteConnection() *SQLiteConnection {
	return &SQLiteConnection{
		mux: &sync.RWMutex{},
	}
}

func (s *SQLiteConnection) Connect(DBConfig DBConfig) error {
	if DBConfig.DBType != SQLite {
		log.Printf("DBType is not SQLite")
		return nil

	}
	conn, err := sql.Open("sqlite3", DBConfig.FilePath)
	if err != nil {
		log.Printf("Failed to connect to sqlite: %v", err)
		return err
	}
	s.conn = conn
	return nil
}

func (s *SQLiteConnection) CreateTablesAndStatements() error {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, err := s.conn.Exec("CREATE TABLE IF NOT EXISTS short_urls (url TEXT PRIMARY KEY, full_url TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, expires_at DATETIME)")
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}
	PreparedStatement, err := s.conn.Prepare("INSERT INTO short_urls (url, full_url, expires_at) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return err
	}
	s.insertNewShortUrl = PreparedStatement
	return nil
}

func (s *SQLiteConnection) InsertNewShortUrl(url, fullUrl string, expiresAt *time.Time) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, err := s.insertNewShortUrl.Exec(url, fullUrl, expiresAt)
	if err != nil {
		log.Printf("Failed to insert new short url: %v", err)
		return err
	}
	return nil
}

func (s *SQLiteConnection) GetFullUrl(url string) (string, *time.Time, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	var fullUrl string
	var expiresAt *time.Time
	err := s.conn.QueryRow("SELECT full_url, expires_at FROM short_urls WHERE url = ?", url).Scan(&fullUrl, &expiresAt)
	if err != nil {
		log.Printf("Failed to get full url: %v", err)
		return "", nil, err
	}
	return fullUrl, expiresAt, nil
}

func (s *SQLiteConnection) Close() {
	s.conn.Close()
}

func (s *SQLiteConnection) Cleanup() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.conn.Exec("DROP TABLE short_urls")
}

func (s *SQLiteConnection) printAllTableContents() {
	s.mux.RLock()
	defer s.mux.RUnlock()
	rows, err := s.conn.Query("SELECT * FROM short_urls")
	if err != nil {
		log.Printf("Failed to print all table contents: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var url, fullUrl string
		var createdAt, expiresAt *time.Time
		err = rows.Scan(&url, &fullUrl, &createdAt, &expiresAt)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
		}
		log.Printf("url: %s, full_url: %s, created_at: %s, expires_at: %s", url, fullUrl, createdAt, expiresAt)
	}
}
