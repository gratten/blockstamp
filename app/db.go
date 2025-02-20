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
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// // connStr := fmt.Sprintf("postgres://postgres:%s@localhost:5432/blockstamp?sslmode=disable", os.Getenv("DB_PASSWORD"))
	// connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	// 	user, password, host, port, dbname)

	// host := getEnvOrDefault("DB_HOST", "postgres")
	// port := getEnvOrDefault("DB_PORT", "5432")
	// user := getEnvOrDefault("DB_USER", "postgres")
	// password := getEnvOrDefault("DB_PASSWORD", "hello")
	// dbname := getEnvOrDefault("DB_NAME", "blockstamp")

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Build URL format for migrations
	migrateConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	// Log the connection string (remove sensitive info for production)
	log.Printf("Connecting to database with host=%s port=%s user=%s dbname=%s",
		host, port, user, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
		return nil, err
	}
	// Run migrations
	migrationsPath := "file://app/db/migrations"
	m, err := migrate.New(migrationsPath, migrateConnStr)
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

// Helper function to get environment variable with default fallback
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
