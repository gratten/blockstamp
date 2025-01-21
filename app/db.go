package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() (*sql.DB, error) {
	// connStr := "user=postgres password=" + os.Getenv("DB_PASSWORD") + " dbname=blockstamp sslmode=disable"
	connStr := fmt.Sprintf("postgres://postgres:%s@localhost:5432/blockstamp?sslmode=disable", os.Getenv("DB_PASSWORD"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
		return nil, err
	}
	// Run migrations
	migrationsPath := "file://app/db/migrations"
	m, err := migrate.New(migrationsPath, connStr)
	if err != nil {
		log.Fatal("Could not create migrate instance:", err)
		return nil, err
	}

	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration failed:", err)
		return nil, err
	}

	log.Println("Migrations applied successfully!")
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
