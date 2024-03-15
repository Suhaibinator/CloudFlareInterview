package dpadapters

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	conn              *sql.DB
	insertNewShortUrl *sql.Stmt
}

func NewPostgresConnection() *PostgresConnection {
	return &PostgresConnection{}
}

func (p *PostgresConnection) Connect(host string, port int, username, password, dbName string) error {

	connStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to postgres: %v", err)
		return err
	}
	p.conn = conn
	return nil
}

func (p *PostgresConnection) CreateTablesAndStatements() error {
	_, err := p.conn.Exec("CREATE TABLE IF NOT EXISTS short_urls (url VARCHAR(255) PRIMARY KEY, full_url VARCHAR(255), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, expires_at TIMESTAMP)")
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}
	PreparedStatement, err := p.conn.Prepare("INSERT INTO short_urls (url, full_url, expires_at) VALUES ($1, $2, $3)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return err
	}
	p.insertNewShortUrl = PreparedStatement
	return nil
}

func (p *PostgresConnection) InsertNewShortUrl(url, fullUrl string, expiresAt *time.Time) error {

	_, err := p.insertNewShortUrl.Exec(url, fullUrl, expiresAt)
	if err != nil {
		log.Printf("Failed to insert new short url: %v", err)
		return err
	}
	return nil
}

func (p *PostgresConnection) GetFullUrl(url string) (string, *time.Time, error) {
	var fullUrl string
	var expiresAt *time.Time
	err := p.conn.QueryRow("SELECT full_url, expires_at FROM short_urls WHERE url = $1", url).Scan(&fullUrl, &expiresAt)
	if err != nil {
		log.Printf("Failed to get full url: %v", err)
		return "", nil, err
	}
	return fullUrl, expiresAt, nil
}

func (p *PostgresConnection) Close() {
	p.conn.Close()
}

func (p *PostgresConnection) Cleanup() {
	p.conn.Exec("DROP TABLE short_urls")
}

func (p *PostgresConnection) PrintAllTableContents() {
	rows, err := p.conn.Query("SELECT * FROM short_urls")
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
