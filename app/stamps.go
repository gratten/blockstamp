// stamps.go
package main

import (
	// "database/sql"
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
	rows, err := db.Query("SELECT blockheight, stamp FROM stamps ORDER BY blockheight DESC")
	if err != nil {
		log.Printf("Error fetching stamps: %v", err)
		http.Error(w, "Failed to fetch stamps", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stamps []struct {
		Blockheight int
		Stamp       string
	}

	// Loop through the result set
	for rows.Next() {
		var stamp struct {
			Blockheight int
			Stamp       string
		}
		err := rows.Scan(&stamp.Blockheight, &stamp.Stamp)
		if err != nil {
			log.Printf("Error scanning stamp: %v", err)
			http.Error(w, "Failed to read stamps", http.StatusInternalServerError)
			return
		}
		stamps = append(stamps, stamp)
	}

	// Render the page with the stamps data
	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.New("stamps").Parse(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <title>Stamps</title>
        </head>
        <body>
            <h1>Stamps</h1>
            <table>
                <tr>
                    <th>Blockheight</th>
                    <th>Stamp</th>
                </tr>
                {{range .}}
                    <tr>
                        <td>{{.Blockheight}}</td>
                        <td>{{.Stamp}}</td>
                    </tr>
                {{end}}
            </table>
        </body>
        </html>
    `)

	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, stamps)
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
