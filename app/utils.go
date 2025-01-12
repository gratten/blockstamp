// utils.go
package main

import (
	"log"
	"os"
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
	// db         *sql.DB
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
