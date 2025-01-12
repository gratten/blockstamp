package main

import (
	// "db"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	// "database/sql"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/joho/godotenv"
	// _ "github.com/lib/pq"
)

// Cache TTL of 24hrs
const cacheTTL = 24 * time.Hour
const averageBlockTime = 600

type cacheEntry struct {
	height    int64
	timestamp time.Time
}

var (
	client     *rpcclient.Client
	once       sync.Once
	cacheMap   = make(map[int64]cacheEntry)
	cacheMutex sync.RWMutex
	// db         *sql.DB
)

// func initDB() (*sql.DB, error) {
// 	// Get the password from environment variables
// 	password := os.Getenv("DB_PASSWORD")
// 	if password == "" {
// 		return nil, fmt.Errorf("database password not set in environment variables")
// 	}

// 	// Build the connection string using the environment variable
// 	connStr := fmt.Sprintf("user=postgres password=%s dbname=blockstamp sslmode=disable", password)
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Optional: Ping the database to ensure connection
// 	if err = db.Ping(); err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }

// func closeDB(db *sql.DB) {
// 	if db != nil {
// 		err := db.Close()
// 		if err != nil {
// 			log.Printf("Error closing the database: %v", err)
// 		}
// 	}
// }

func init() {
	once.Do(func() {
		// Load .env file once
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}

		username := os.Getenv("BTCUSER")
		password := os.Getenv("PASSWORD")
		host := os.Getenv("HOST")
		log.Println("host: ", host)

		// Initialize rpcclient only once
		connCfg := &rpcclient.ConnConfig{
			Host:         host,
			User:         username,
			Pass:         password,
			HTTPPostMode: true,
			DisableTLS:   true,
		}

		var err2 error
		client, err2 = rpcclient.New(connCfg, nil)
		if err2 != nil {
			log.Fatal(err2)
		}
	})
}

func clearCachePeriodically() {
	for {
		time.Sleep(10 * time.Minute)

		cacheMutex.Lock()
		for targetTime, entry := range cacheMap {
			if time.Since(entry.timestamp) > cacheTTL {
				delete(cacheMap, targetTime)
			}
		}
		cacheMutex.Unlock()
	}
}

func getBlockTime(height int64) (int64, error) {
	hash, err := client.GetBlockHash(height)
	if err != nil {
		return 0, fmt.Errorf("error getting block hash: %v", err)
	}
	block, err := client.GetBlockVerbose(hash)
	if err != nil {
		return 0, fmt.Errorf("error getting block details: %v", err)
	}
	return block.Time, nil
}

func binarySearch(blockCount int64, targetTime int64) string {
	// Check if the result is already cached
	cacheMutex.RLock()
	if cachedEntry, ok := cacheMap[targetTime]; ok {
		cacheMutex.RUnlock()
		resultStr := strconv.FormatInt(cachedEntry.height, 10)
		return resultStr
	}
	cacheMutex.RUnlock()

	// Get the latest block's time
	latestBlockHash, err := client.GetBlockHash(blockCount)
	if err != nil {
		log.Fatal(err)
	}

	latestBlock, err := client.GetBlockVerbose(latestBlockHash)
	if err != nil {
		log.Fatal(err)
	}

	latestBlockTime := latestBlock.Time

	// Check if the target time is in the future
	if targetTime > latestBlockTime {
		timeDifference := targetTime - latestBlockTime
		estimatedFutureBlocks := timeDifference / averageBlockTime
		resultStr := strconv.FormatInt(blockCount+estimatedFutureBlocks, 10) + " (estimate)"
		return resultStr
	}

	// Perform binary search for past blocks
	var leftBlockHeight, rightBlockHeight int64 = 0, blockCount

	for leftBlockHeight <= rightBlockHeight {
		midBlockHeight := (leftBlockHeight + rightBlockHeight) / 2

		midBlockTime, err := getBlockTime(midBlockHeight)
		if err != nil {
			log.Printf("Error getting block time: %v", err)
			return "error"
		}

		if midBlockTime == targetTime {
			cacheMutex.Lock()
			cacheMap[targetTime] = cacheEntry{height: midBlockHeight, timestamp: time.Now()}
			cacheMutex.Unlock()
			resultStr := strconv.FormatInt(midBlockHeight, 10)
			return resultStr
		} else if midBlockTime < targetTime {
			leftBlockHeight = midBlockHeight + 1
		} else {
			rightBlockHeight = midBlockHeight - 1
		}
	}

	result := leftBlockHeight
	cacheMutex.Lock()
	cacheMap[targetTime] = cacheEntry{height: result, timestamp: time.Now()}
	cacheMutex.Unlock()
	resultStr := strconv.FormatInt(result, 10)
	return resultStr
}

// // Handler to show all stamps
// func showStamps(w http.ResponseWriter, r *http.Request) {
// 	// Query the database for all stamps
// 	rows, err := db.Query("SELECT blockheight, stamp FROM stamps ORDER BY blockheight DESC")
// 	if err != nil {
// 		log.Printf("Error fetching stamps: %v", err)
// 		http.Error(w, "Failed to fetch stamps", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var stamps []struct {
// 		Blockheight int
// 		Stamp       string
// 	}

// 	// Loop through the result set
// 	for rows.Next() {
// 		var stamp struct {
// 			Blockheight int
// 			Stamp       string
// 		}
// 		err := rows.Scan(&stamp.Blockheight, &stamp.Stamp)
// 		if err != nil {
// 			log.Printf("Error scanning stamp: %v", err)
// 			http.Error(w, "Failed to read stamps", http.StatusInternalServerError)
// 			return
// 		}
// 		stamps = append(stamps, stamp)
// 	}

// 	// Render the page with the stamps data
// 	w.Header().Set("Content-Type", "text/html")
// 	tmpl, err := template.New("stamps").Parse(`
// 		<!DOCTYPE html>
// 		<html lang="en">
// 		<head>
// 			<meta charset="UTF-8">
// 			<title>Stamps</title>
// 		</head>
// 		<body>
// 			<h1>Stamps</h1>
// 			<table>
// 				<tr>
// 					<th>Blockheight</th>
// 					<th>Stamp</th>
// 				</tr>
// 				{{range .}}
// 					<tr>
// 						<td>{{.Blockheight}}</td>
// 						<td>{{.Stamp}}</td>
// 					</tr>
// 				{{end}}
// 			</table>
// 		</body>
// 		</html>
// 	`)

// 	if err != nil {
// 		log.Printf("Error parsing template: %v", err)
// 		http.Error(w, "Failed to render page", http.StatusInternalServerError)
// 		return
// 	}

// 	err = tmpl.Execute(w, stamps)
// 	if err != nil {
// 		log.Printf("Error executing template: %v", err)
// 		http.Error(w, "Failed to render page", http.StatusInternalServerError)
// 		return
// 	}
// }

func main() {
	defer client.Shutdown()

	// initialize db
	var err error
	db, err := initDB()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer closeDB(db) // Ensure the database connection is closed when the program exits

	// Start cache cleanup goroutine
	go clearCachePeriodically()

	fmt.Println("Server started")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		blockheight := "Enter a date to find the blockheight."
		tmpl.Execute(w, blockheight)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		blockCount, err := client.GetBlockCount()
		if err != nil {
			log.Fatal(err)
		}

		year, _ := strconv.Atoi(r.PostFormValue("year"))
		month, _ := strconv.Atoi(r.PostFormValue("month"))
		day, _ := strconv.Atoi(r.PostFormValue("day"))
		hour, _ := strconv.Atoi(r.PostFormValue("hour"))
		minute, _ := strconv.Atoi(r.PostFormValue("minute"))
		second, _ := strconv.Atoi(r.PostFormValue("second"))

		location, err := time.LoadLocation("America/New_York")
		if err != nil {
			fmt.Println("Error loading location:", err)
		}

		givenDateTime := time.Date(year, time.Month(month), day, hour, minute, second, 0, location)
		targetTime := givenDateTime.Unix()

		resultStr := binarySearch(blockCount, targetTime)

		response := map[string]interface{}{
			"blockheight":   resultStr,
			"showStampForm": strings.Contains(resultStr, " (estimate)"), // Flag for the stamp form
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		duration := time.Since(start)
		log.Printf("Time taken for request: %v", duration)
	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		blockCount, err := client.GetBlockCount()
		if err != nil {
			http.Error(w, "Unable to fetch block count", http.StatusInternalServerError)
			return
		}

		blockHash, err := client.GetBlockHash(blockCount)
		if err != nil {
			http.Error(w, "Unable to fetch block hash", http.StatusInternalServerError)
			return
		}

		block, err := client.GetBlockVerbose(blockHash)
		if err != nil {
			http.Error(w, "Unable to fetch block details", http.StatusInternalServerError)
			return
		}

		blockTime := time.Unix(block.Time, 0).Format("2006-01-02 15:04:05")

		response := fmt.Sprintf("Current Blockheight: %d<br>Mined on: %s", blockCount, blockTime)

		w.Header().Set("Content-Type", "text/html") // Return HTML since we include <br> tags
		fmt.Fprint(w, response)
	}
	h4 := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm() // Parse form data
		blockheightStr := r.PostFormValue("blockheight")
		re := regexp.MustCompile("[^0-9]")
		blockheightStr = re.ReplaceAllString(blockheightStr, "")
		log.Println("blockheightStr: ", blockheightStr)
		stamp := r.PostFormValue("stamp")
		log.Println("stamp: ", stamp)

		// Convert blockheight string to integer
		blockheight, err := strconv.Atoi(blockheightStr)
		if err != nil {
			log.Printf("Error converting blockheight: %v", err)
			http.Error(w, "Invalid blockheight", http.StatusBadRequest)
			return
		}

		// Insert into the database
		insertStmt := `INSERT INTO stamps (blockheight, stamp) VALUES ($1, $2)`
		_, err = db.Exec(insertStmt, blockheight, stamp)
		if err != nil {
			log.Printf("Error inserting stamp: %v", err)
			http.Error(w, "Failed to insert stamp", http.StatusInternalServerError)
			return
		}

		// Respond with success (optional)
		fmt.Fprintf(w, "Stamp submitted successfully!")
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/get-blockheight/", h2)
	http.HandleFunc("/current-blockheight/", h3)
	http.HandleFunc("/submit-stamp/", h4)
	// http.HandleFunc("/stamps", showStamps) // Add this line for the /stamps route
	log.Fatal(http.ListenAndServe(":8000", nil))
}
