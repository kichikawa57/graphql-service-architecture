package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "graphql_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	maxRetries := 5
	retryDelay := time.Second * 2

	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Attempt %d failed to open database: %v", i+1, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				retryDelay *= 2
				continue
			}
			return nil, fmt.Errorf("failed to open database after %d attempts: %w", maxRetries, err)
		}

		if err := db.Ping(); err != nil {
			log.Printf("Attempt %d failed to ping database: %v", i+1, err)
			db.Close()
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				retryDelay *= 2
				continue
			}
			return nil, fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, err)
		}

		log.Println("Database connection established")
		return &DB{db}, nil
	}

	return nil, fmt.Errorf("failed to establish database connection after %d attempts", maxRetries)
}

func (db *DB) CreateTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Users table created successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
