// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"fmt"
	"time"
    "os"
	"github.com/btcsuite/btcd/rpcclient"
    "github.com/joho/godotenv"
)

func binarySearch(client *rpcclient.Client, blockCount int64, targetTime int64) int64 {
	var leftBlockHeight, rightBlockHeight int64 = 0, blockCount

	for leftBlockHeight <= rightBlockHeight {

		// fmt.Println(rightBlockHeight)
		// fmt.Println(blockCount)
		// fmt.Println(targetTime)

		midBlockHeight := (leftBlockHeight + rightBlockHeight) / 2

		midBlockHash, err := client.GetBlockHash(midBlockHeight)
			if err != nil {
				log.Fatal(err)
			}

		midBlock, err := client.GetBlockVerbose(midBlockHash)
			if err != nil {
				log.Fatal(err)
			}
		// fmt.Println(midBlock.Time)

		// leftBlockHash, err := client.GetBlockHash(leftBlockHeight)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// leftBlock, err := client.GetBlockVerbose(leftBlockHash)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// rightBlockHash, err := client.GetBlockHash(rightBlockHeight)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// rightBlock, err := client.GetBlockVerbose(rightBlockHash)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		midBlockTime := midBlock.Time
		// leftBlockTime := leftBlock.Time
		// rightBlockTime := rightBlock.Time

		if midBlockTime == targetTime {
			return midBlockHeight
		} else if midBlockTime < targetTime {
			leftBlockHeight = midBlockHeight + 1
		} else {
			rightBlockHeight = midBlockHeight - 1
		}
	}
	return(leftBlockHeight)

}

func main() {
    // Load the environment variables from the .env file
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
    username := os.Getenv("BTCUSER")
    password := os.Getenv("PASSWORD")
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig {
		Host:         "localhost:8332",
		User:         username,
		Pass:         password,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("Block count: %d", blockCount)

	var year, month, day, hour, minute, sec int
	fmt.Println("Note: assuming eastern time")
	fmt.Println("Enter year:")
	fmt.Scan(&year)
	fmt.Println("Enter month:")
	fmt.Scan(&month)
	fmt.Println("Enter day:")
	fmt.Scan(&day)
	fmt.Println("Enter hour:")
	fmt.Scan(&hour)
	fmt.Println("Enter minute:")
	fmt.Scan(&minute)
	fmt.Println("Enter second:")
	fmt.Scan(&sec)

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		// Handle error, e.g., log it or return an error
		fmt.Println("Error loading location:", err)
		// return or handle the error as appropriate
	}

	givenDateTime := time.Date(year, time.Month(month), day, hour, minute, sec, 0, location)
	// fmt.Printf("given date time: %d\n", givenDateTime)
	// fmt.Printf("given date time: %d\n", givenDateTime.UTC())
	targetTime := givenDateTime.Unix()
	// fmt.Printf("target time: %d\n", targetTime)


	fmt.Println("Finding block height...")
	result := binarySearch(client, blockCount, targetTime)
	// fmt.Println(result)
	fmt.Printf("The block height at this date and time was: %d\n", result)
}
