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
	"strings"
)

func ShowStampsTable(w http.ResponseWriter, r *http.Request) {
	// log.Printf("request from /stamps-table")

	// Query the database for all stamps
	rows, err := db.Query("SELECT blockheight, stamp, txid FROM stamps ORDER BY blockheight DESC")
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
		TxIDColor   string
	}

	// Loop through the result set
	for rows.Next() {
		var stamp struct {
			Blockheight int
			Stamp       string
			TxID        sql.NullString
			TxIDColor   string
		}
		err := rows.Scan(&stamp.Blockheight, &stamp.Stamp, &stamp.TxID)
		if err != nil {
			log.Printf("Error scanning stamp: %v", err)
			http.Error(w, "Failed to read stamps", http.StatusInternalServerError)
			return
		}

		// Determine if the transaction is mined
		if stamp.TxID.Valid {
			mined, err := CheckIfTransactionMined(stamp.TxID.String)
			if err != nil {
				// log.Printf("Error checking transaction status: %v", err)
				stamp.TxIDColor = "black"
			} else if mined {
				stamp.TxIDColor = "green"
			} else {
				stamp.TxIDColor = "black"
			}
		} else {
			stamp.TxIDColor = "black"
		}

		stamps = append(stamps, stamp)
	}

	// Render only the rows of the table
	tmpl, err := template.ParseFiles("app/stamps-table.html")
	if err != nil {
		log.Printf("Error parsing table template: %v", err)
		http.Error(w, "Failed to render table rows", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, stamps)
	if err != nil {
		log.Printf("Error executing table template: %v", err)
		http.Error(w, "Failed to render table rows", http.StatusInternalServerError)
		return
	}
}

func ShowStamps(w http.ResponseWriter, r *http.Request) {
	log.Printf("request from /stamps")

	// Declare the tmpl variable
	tmpl, err := template.ParseFiles("app/layout.html", "app/stamps.html", "app/blockheight.html")
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Failed to load templates", http.StatusInternalServerError)
		return
	}

	// Query the database for all stamps
	rows, err := db.Query("SELECT blockheight, stamp, txid FROM stamps ORDER BY blockheight DESC")
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
		TxIDColor   string
	}

	// Loop through the result set and populate the stamps array
	for rows.Next() {
		var stamp struct {
			Blockheight int
			Stamp       string
			TxID        sql.NullString
			TxIDColor   string
		}
		err := rows.Scan(&stamp.Blockheight, &stamp.Stamp, &stamp.TxID)
		if err != nil {
			log.Printf("Error scanning stamp: %v", err)
			http.Error(w, "Failed to read stamps", http.StatusInternalServerError)
			return
		}

		// Check if the transaction is mined and set color
		if stamp.TxID.Valid {
			mined, err := CheckIfTransactionMined(stamp.TxID.String)
			if err != nil {
				// log.Printf("Error checking transaction status: %v", err)
				stamp.TxIDColor = "black"
			} else if mined {
				stamp.TxIDColor = "green"
			} else {
				stamp.TxIDColor = "black"
			}
		} else {
			stamp.TxIDColor = "black" // Handle NULL TxID
		}

		stamps = append(stamps, stamp)
	}

	// Check if we are making an HTMX request (this part is key to handling partial updates)
	isHTMX := strings.Contains(r.Header.Get("HX-Request"), "true")

	if isHTMX {
		// Render only the table rows as HTML
		err = tmpl.ExecuteTemplate(w, "stamps.html", struct {
			Stamps []struct {
				Blockheight int
				Stamp       string
				TxID        sql.NullString
				TxIDColor   string
			}
		}{Stamps: stamps})
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
			return
		}
	} else {
		// If it's not an HTMX request, render the full page
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

	// // txid, err := transaction(blockheight, stamp)
	// txid, err := transaction(blockheight, stamp)
	// if err != nil {
	// 	// Handle the error
	// 	log.Printf("Error obtaining txid: %v", err)
	// }

	// // Insert into the database
	// insertStmt := `INSERT INTO stamps (blockheight, stamp, txid) VALUES ($1, $2, $3)`
	// _, err = db.Exec(insertStmt, blockheight, stamp, txid)
	// if err != nil {
	// 	log.Printf("Error inserting stamp: %v", err)
	// 	http.Error(w, "Failed to insert stamp", http.StatusInternalServerError)
	// 	return
	// }

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
	// transaction(blockheight, stamp)
}
