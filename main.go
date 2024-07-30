package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
	return (leftBlockHeight)
}

func main() {
	// blockchain stuff
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	username := os.Getenv("BTCUSER")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		// Host:         "localhost:8332",
		Host:         host,
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

	// // Get the current block count.
	// blockCount, err := client.GetBlockCount()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// web stuff (collect target time)
	fmt.Println("hello world")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		blockheight := "Enter a date to find the blockheight."
		tmpl.Execute(w, blockheight)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {

		// Get the current block count.
		blockCount, err := client.GetBlockCount()
		if err != nil {
			log.Fatal(err)
		}
		// log.Print("HTMX request recieved")
		// log.Print(r.Header.Get("HX-Request"))
		year, _ := strconv.Atoi(r.PostFormValue("year"))
		month, _ := strconv.Atoi(r.PostFormValue("month"))
		day, _ := strconv.Atoi(r.PostFormValue("day"))
		hour, _ := strconv.Atoi(r.PostFormValue("hour"))
		minute, _ := strconv.Atoi(r.PostFormValue("minute"))
		second, _ := strconv.Atoi(r.PostFormValue("second"))

		location, err := time.LoadLocation("America/New_York")
		if err != nil {
			// Handle error, e.g., log it or return an error
			fmt.Println("Error loading location:", err)
			// return or handle the error as appropriate
		}

		givenDateTime := time.Date(year, time.Month(month), day, hour, minute, second, 0, location)
		// // fmt.Printf("given date time: %d\n", givenDateTime)
		// // fmt.Printf("given date time: %d\n", givenDateTime.UTC())
		targetTime := givenDateTime.Unix()
		// log.Print(targetTime)
		// target := FormatInt(int64(targetTime), 10)
		// // fmt.Println(year)

		// fmt.Println("Finding block height...")
		result := binarySearch(client, blockCount, targetTime)
		// fmt.Println(result)
		// fmt.Printf("The block height at this date and time was: %d\n", result)
		// fmt.Printf("Type of myVar: %T\n", result)
		resultStr := strconv.FormatInt(result, 10)
		tmpl, _ := template.New("t").Parse(resultStr)
		tmpl.Execute(w, nil)

	}
	http.HandleFunc("/", h1)
	http.HandleFunc("/get-blockheight/", h2)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
