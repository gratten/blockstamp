package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	connStr := "user=postgres password=" + os.Getenv("DB_PASSWORD") + " dbname=blockstamp sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
		return nil, err
	}
	return db, nil
}

func closeDB(db *sql.DB) {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing the database: %v", err)
		}
	}
}
