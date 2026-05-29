package db

import (
	"database/sql"
	"log"
)

func InitiateDatabaseConnection() (*sql.DB, error) {
	log.Println("Starting Database Connection...")
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	return sql.Open("postgres", connStr)
}

func ShutdownDatabaseConnection(DB *sql.DB) {
	log.Println("Database Connection shutdown...")
	DB.Close()
}
