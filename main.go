package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/joho/godotenv"
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
)

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
		log.Printf("future")
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

func main() {
	defer client.Shutdown()

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

		duration := time.Since(start)
		log.Printf("Time taken for request: %v", duration)

		tmpl, _ := template.New("t").Parse(resultStr)
		tmpl.Execute(w, nil)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/get-blockheight/", h2)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
