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

func newPostgresConnection() *PostgresConnection {
	return &PostgresConnection{}
}

func (p *PostgresConnection) Connect(DBConfig DBConfig) error {

	if DBConfig.DBType != Postgres {
		log.Printf("DBType is not Postgres")
		return nil
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBConfig.Host, DBConfig.Port, DBConfig.Username, DBConfig.Password, DBConfig.DBName)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to postgres: %v", err)
		return err
	}
	p.conn = conn
	return nil
}

func (p *PostgresConnection) CreateTablesAndStatements(workerId string) error {
	_, err := p.conn.Exec("CREATE TABLE IF NOT EXISTS short_urls (url VARCHAR(255) PRIMARY KEY, full_url VARCHAR(255), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, expires_at TIMESTAMP)")
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}

	_, err = p.conn.Exec(`
    CREATE TABLE IF NOT EXISTS last_count (worker_id TEXT, count INT);
    INSERT INTO last_count (worker_id, count)
    SELECT ?, 0 WHERE NOT EXISTS (SELECT 1 FROM last_count WHERE worker_id = ?);`, workerId, workerId)
	if err != nil {
		log.Printf("Failed to create counter table: %v", err)
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

func (p *PostgresConnection) DeleteShortUrl(url string) error {
	_, err := p.conn.Exec("DELETE FROM short_urls WHERE url = $1", url)
	if err != nil {
		log.Printf("Failed to delete short url: %v", err)
	}
	return err
}

func (p *PostgresConnection) UpdateCounter(workerID string, newCount int) error {
	_, err := p.conn.Exec("UPDATE last_count SET count = $1 WHERE worker_id = $2", newCount, workerID)
	if err != nil {
		log.Printf("Failed to update counter: %v", err)
	}
	return err
}

func (p *PostgresConnection) GetCounter(workerID string) (int, error) {
	var count int
	err := p.conn.QueryRow("SELECT count FROM last_count WHERE worker_id = $1", workerID).Scan(&count)
	if err != nil {
		log.Printf("Failed to get counter: %v", err)
	}
	return count, err
}

func (p *PostgresConnection) Close() {
	p.conn.Close()
}

func (p *PostgresConnection) Cleanup() {
	p.conn.Exec("DROP TABLE short_urls")
}

func (p *PostgresConnection) printAllTableContents() {
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
