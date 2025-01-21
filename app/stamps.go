// stamps.go
package main

import (
	// "database/sql"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// ShowStamps handler to show all stamps
func ShowStamps(w http.ResponseWriter, r *http.Request) {
	// Query the database for all stamps
	// log.Println("Database connection:", db) // Logs the db connection (should not be nil)
	// rows, err := db.Query("SELECT blockheight, stamp FROM stamps ORDER BY id ASC")
	rows, err := db.Query("SELECT blockheight, stamp, txid FROM stamps ORDER BY id ASC")
	if err != nil {
		log.Printf("Error fetching stamps: %v", err)
		http.Error(w, "Failed to fetch stamps", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stamps []struct {
		Blockheight int
		Stamp       string
		TxID        sql.NullString
	}

	// Loop through the result set
	for rows.Next() {
		var stamp struct {
			Blockheight int
			Stamp       string
			TxID        sql.NullString
		}
		err := rows.Scan(&stamp.Blockheight, &stamp.Stamp, &stamp.TxID)
		if err != nil {
			log.Printf("Error scanning stamp: %v", err)
			http.Error(w, "Failed to read stamps", http.StatusInternalServerError)
			return
		}
		stamps = append(stamps, stamp)
	}

	// Render the page with the stamps data
	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.ParseFiles(
		"app/layout.html",
		"app/stamps.html",
		"app/blockheight.html",
	)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}

	// Pass data for layout and stamps table
	err = tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Title":  "Stamps",
		"Stamps": stamps,
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}

// SubmitStamp handler to submit a stamp
func SubmitStamp(w http.ResponseWriter, r *http.Request) {
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

	// txid, err := transaction(blockheight, stamp)
	txid, err := transaction(blockheight, stamp)
	if err != nil {
		// Handle the error
		log.Printf("Error obtaining txid: %v", err)
	}

	// Insert into the database
	insertStmt := `INSERT INTO stamps (blockheight, stamp, txid) VALUES ($1, $2, $3)`
	_, err = db.Exec(insertStmt, blockheight, stamp, txid)
	if err != nil {
		log.Printf("Error inserting stamp: %v", err)
		http.Error(w, "Failed to insert stamp", http.StatusInternalServerError)
		return
	}

	// Respond with success (optional)
	fmt.Fprintf(w, "Stamp submitted successfully!")
	// transaction(blockheight, stamp)
}
