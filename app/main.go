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
	// log.Println(currentBlockHeight)

	// Query the database for rows where target_blockheight <= currentBlockHeight and txid IS NULL
	rows, err := db.Query(`
		SELECT id, blockheight, stamp 
		FROM stamps 
		WHERE blockheight = $1 AND txid IS NULL
	`, currentBlockHeight)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		return
	}
	defer rows.Close()

	// Iterate over the eligible rows
	for rows.Next() {
		var id int
		var targetBlockHeight int
		var opReturnMessage string

		// Scan the row into variables
		if err := rows.Scan(&id, &targetBlockHeight, &opReturnMessage); err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}

		// Call the transaction function
		txID, err := transaction(targetBlockHeight, opReturnMessage)
		if err != nil {
			log.Printf("Error broadcasting transaction for ID %d: %v\n", id, err)
			continue
		}

		// Update the database with the transaction ID
		_, err = db.Exec(`
			UPDATE stamps 
			SET txid = $1
			WHERE id = $2
		`, txID, id)
		if err != nil {
			log.Printf("Error updating transaction ID for ID %d: %v\n", id, err)
		} else {
			log.Printf("Transaction broadcasted successfully for ID %d. TXID: %s\n", id, txID)
		}
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v\n", err)
	}
}

// startChecker periodically checks the database for transactions to broadcast.
func startChecker(db *sql.DB, interval time.Duration) {
	for {
		checkAndBroadcastTransactions(db)
		time.Sleep(interval)
	}
}

var pendingStamps *PendingStamps

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

	pendingStamps = NewPendingStamps()
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/get-blockheight/", GetBlockheightByDate)
	http.HandleFunc("/current-blockheight/", GetCurrentBlockheight)
	http.HandleFunc("/submit-stamp/", SubmitStamp)
	http.HandleFunc("/stamps", ShowStamps)
	http.HandleFunc("/stamps-table", ShowStampsTable)
	http.HandleFunc("/check-payment-status", CheckPaymentStatus)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
