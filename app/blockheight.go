// blockheight.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
)

// GetBlockTime retrieves the block time for a given block height
func GetBlockTime(height int64) (int64, error) {
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

// BinarySearch performs a binary search to find the block height for a given target time
func BinarySearch(blockCount int64, targetTime int64) string {
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

		midBlockTime, err := GetBlockTime(midBlockHeight)
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

// Handler for /get-blockheight/ to fetch the blockheight for a given date and time
func GetBlockheightByDate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}

	// Parse date/time values from the request
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

	// Combine parsed values into a time object
	givenDateTime := time.Date(year, time.Month(month), day, hour, minute, second, 0, location)
	targetTime := givenDateTime.Unix()

	// Use binary search to find the blockheight for the target time
	resultStr := BinarySearch(blockCount, targetTime)

	// Send response
	response := map[string]interface{}{
		"blockheight":   resultStr,
		"showStampForm": strings.Contains(resultStr, " (estimate)"), // Flag for the stamp form
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	duration := time.Since(start)
	log.Printf("Time taken for request: %v", duration)
}

// // Handler for /current-blockheight/ to show the current blockheight
// func GetCurrentBlockheight(w http.ResponseWriter, r *http.Request) {
// 	blockCount, err := client.GetBlockCount()
// 	if err != nil {
// 		http.Error(w, "Unable to fetch block count", http.StatusInternalServerError)
// 		return
// 	}

// 	blockHash, err := client.GetBlockHash(blockCount)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch block hash", http.StatusInternalServerError)
// 		return
// 	}

// 	block, err := client.GetBlockVerbose(blockHash)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch block details", http.StatusInternalServerError)
// 		return
// 	}

// 	blockTime := time.Unix(block.Time, 0).Format("2006-01-02 15:04:05")

// 	// Response format
// 	response := fmt.Sprintf("Current Blockheight: %d<br>Mined on: %s", blockCount, blockTime)

// 	w.Header().Set("Content-Type", "text/html") // Return HTML since we include <br> tags
// 	fmt.Fprint(w, response)
// }

func GetCurrentBlockHeight(client *rpcclient.Client) (int, error) {
	blockCount, err := client.GetBlockCount()
	if err != nil {
		return 0, fmt.Errorf("unable to fetch block count: %w", err)
	}
	return int(blockCount), nil
}

func GetCurrentBlockheight(w http.ResponseWriter, r *http.Request) {
	log.Printf("request from /current-blockheight")
	blockCount, err := GetCurrentBlockHeight(client) // Reuse core function
	if err != nil {
		http.Error(w, "Unable to fetch block count", http.StatusInternalServerError)
		return
	}

	blockHash, err := client.GetBlockHash(int64(blockCount))
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

	// Response format
	response := fmt.Sprintf("Current Blockheight: %d<br>Mined on: %s", blockCount, blockTime)

	w.Header().Set("Content-Type", "text/html") // Return HTML since we include <br> tags
	fmt.Fprint(w, response)
}
