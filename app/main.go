package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

// checkAndBroadcastTransactions checks the database for eligible transactions and broadcasts them.
func checkAndBroadcastTransactions(db *sql.DB) {
	// Retrieve the current block height
	currentBlockHeight, err := GetCurrentBlockHeight(client)
	if err != nil {
		log.Printf("Error getting current block height: %v\n", err)
		return
	}
	log.Println(currentBlockHeight)
}

// startChecker periodically checks the database for transactions to broadcast.
func startChecker(db *sql.DB, interval time.Duration) {
	for {
		checkAndBroadcastTransactions(db)
		time.Sleep(interval)
	}
}

func main() {
	defer client.Shutdown()

	// initialize db
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer closeDB(db) // Ensure the database connection is closed when the program exits

	// Start cache cleanup goroutine
	go clearCachePeriodically()

	fmt.Println("Server started")

	// Start the periodic checker (e.g., every 10 seconds)
	go startChecker(db, 3*time.Second)

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/get-blockheight/", GetBlockheightByDate)
	http.HandleFunc("/current-blockheight/", GetCurrentBlockheight)
	http.HandleFunc("/submit-stamp/", SubmitStamp)
	http.HandleFunc("/stamps", ShowStamps)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
