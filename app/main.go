package main

import (
	"fmt"
	"log"
	"net/http"
)

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

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/get-blockheight/", GetBlockheightByDate)
	http.HandleFunc("/current-blockheight/", GetCurrentBlockheight)
	http.HandleFunc("/submit-stamp/", SubmitStamp)
	http.HandleFunc("/stamps", ShowStamps)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
