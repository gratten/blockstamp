// stamps.go
package main

import (
	// "database/sql"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
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

// // SubmitStamp handler to submit a stamp
// func SubmitStamp(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm() // Parse form data
// 	blockheightStr := r.PostFormValue("blockheight")
// 	re := regexp.MustCompile("[^0-9]")
// 	blockheightStr = re.ReplaceAllString(blockheightStr, "")
// 	log.Println("blockheightStr: ", blockheightStr)
// 	stamp := r.PostFormValue("stamp")
// 	log.Println("stamp: ", stamp)

// 	// Convert blockheight string to integer
// 	blockheight, err := strconv.Atoi(blockheightStr)
// 	if err != nil {
// 		log.Printf("Error converting blockheight: %v", err)
// 		http.Error(w, "Invalid blockheight", http.StatusBadRequest)
// 		return
// 	}

// 	// // txid, err := transaction(blockheight, stamp)
// 	// txid, err := transaction(blockheight, stamp)
// 	// if err != nil {
// 	// 	// Handle the error
// 	// 	log.Printf("Error obtaining txid: %v", err)
// 	// }

// 	// // Insert into the database
// 	// insertStmt := `INSERT INTO stamps (blockheight, stamp, txid) VALUES ($1, $2, $3)`
// 	// _, err = db.Exec(insertStmt, blockheight, stamp, txid)
// 	// if err != nil {
// 	// 	log.Printf("Error inserting stamp: %v", err)
// 	// 	http.Error(w, "Failed to insert stamp", http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// Insert into the database
// 	insertStmt := `INSERT INTO stamps (blockheight, stamp) VALUES ($1, $2)`
// 	_, err = db.Exec(insertStmt, blockheight, stamp)
// 	if err != nil {
// 		log.Printf("Error inserting stamp: %v", err)
// 		http.Error(w, "Failed to insert stamp", http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with success (optional)
// 	fmt.Fprintf(w, "Stamp submitted successfully!")
// 	// transaction(blockheight, stamp)
// }

// func SubmitStamp(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	blockheightStr := r.PostFormValue("blockheight")
// 	stamp := r.PostFormValue("stamp")

// 	// Create invoice for the stamp
// 	invoice, err := createInvoice(fmt.Sprintf("Blockstamp: %s", stamp))
// 	if err != nil {
// 		http.Error(w, "Failed to create invoice", http.StatusInternalServerError)
// 		return
// 	}

// 	// Instead of inserting to database now, store the stamp data temporarily
// 	// We can use a cache or temporary storage for pending payments
// 	pendingStamps.Set(invoice.PaymentHash, StampData{
// 		Blockheight: blockheightStr,
// 		Stamp:       stamp,
// 	})

// 	// Return invoice data to client
// 	response := struct {
// 		PaymentRequest string `json:"payment_request"`
// 		PaymentHash    string `json:"payment_hash"`
// 	}{
// 		PaymentRequest: invoice.PaymentRequest,
// 		PaymentHash:    invoice.PaymentHash,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

func SubmitStamp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	blockheightStr := r.PostFormValue("blockheight")
	stamp := r.PostFormValue("stamp")

	// Create invoice for the stamp
	invoice, err := createInvoice(fmt.Sprintf("Blockstamp: %s", stamp))
	if err != nil {
		http.Error(w, "Failed to create invoice", http.StatusInternalServerError)
		return
	}

	// Instead of inserting to database now, store the stamp data temporarily
	pendingStamps.Set(invoice.PaymentHash, StampData{
		Blockheight: blockheightStr,
		Stamp:       stamp,
	})

	// Start a goroutine to poll payment status
	go func() {
		for i := 0; i < 60; i++ { // Try for 5 minutes
			time.Sleep(5 * time.Second)

			// Check payment status using LNBits API
			url := fmt.Sprintf("%s/api/v1/payments/%s", os.Getenv("LNBITS_URL"), invoice.PaymentHash)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("X-Api-Key", os.Getenv("LNBITS_API_KEY"))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error checking payment status: %v", err)
				continue
			}
			defer resp.Body.Close()

			var payment struct {
				Paid bool `json:"paid"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&payment); err != nil {
				log.Printf("Error decoding payment response: %v", err)
				continue
			}

			if payment.Paid {
				// Get the pending stamp data
				stampData, exists := pendingStamps.Get(invoice.PaymentHash)
				if !exists {
					log.Printf("No pending stamp found for payment hash: %s", invoice.PaymentHash)
					return
				}

				// Clean and convert blockheight string to integer
				blockheightStr := strings.Split(stampData.Blockheight, " ")[0] // Remove " (estimate)" if present
				blockheight, err := strconv.Atoi(blockheightStr)
				if err != nil {
					log.Printf("Error converting blockheight to integer: %v", err)
					return
				}

				// Insert into database
				insertStmt := `INSERT INTO stamps (blockheight, stamp) VALUES ($1, $2)`
				_, err = db.Exec(insertStmt, blockheight, stampData.Stamp)
				if err != nil {
					log.Printf("Error inserting stamp into database: %v", err)
					return
				}

				// Clean up
				pendingStamps.Delete(invoice.PaymentHash)
				log.Printf("Successfully processed payment and stored stamp")
				return
			}
		}
		log.Printf("Payment timeout for hash: %s", invoice.PaymentHash)
	}()

	// Generate QR code
	qr, err := qrcode.Encode(invoice.PaymentRequest, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error generating QR code: %v", err)
	}
	qrBase64 := base64.StdEncoding.EncodeToString(qr)

	// // Return invoice data to client
	// response := struct {
	// 	PaymentRequest string `json:"payment_request"`
	// 	PaymentHash    string `json:"payment_hash"`
	// }{
	// 	PaymentRequest: invoice.PaymentRequest,
	// 	PaymentHash:    invoice.PaymentHash,
	// }

	// Return invoice data to client
	response := struct {
		PaymentRequest string `json:"payment_request"`
		PaymentHash    string `json:"payment_hash"`
		QRCode         string `json:"qr_code"`
	}{
		PaymentRequest: invoice.PaymentRequest,
		PaymentHash:    invoice.PaymentHash,
		QRCode:         qrBase64,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PaymentWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		PaymentHash string `json:"payment_hash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Get the pending stamp data
	stampData, exists := pendingStamps.Get(payload.PaymentHash)
	if !exists {
		http.Error(w, "Invalid payment hash", http.StatusBadRequest)
		return
	}

	// Only now insert into database after payment confirmed
	insertStmt := `INSERT INTO stamps (blockheight, stamp) VALUES ($1, $2)`
	_, err := db.Exec(insertStmt, stampData.Blockheight, stampData.Stamp)
	if err != nil {
		http.Error(w, "Failed to insert stamp", http.StatusInternalServerError)
		return
	}

	// Clean up the pending payment
	pendingStamps.Delete(payload.PaymentHash)

	w.WriteHeader(http.StatusOK)
}

func CheckPaymentStatus(w http.ResponseWriter, r *http.Request) {
	paymentHash := r.URL.Query().Get("payment_hash")

	// Check if payment exists in pendingStamps
	_, exists := pendingStamps.Get(paymentHash)
	if !exists {
		// If it doesn't exist in pending, it was successful
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"paid": true})
		return
	}

	// Check LNBits payment status
	url := fmt.Sprintf("%s/api/v1/payments/%s", os.Getenv("LNBITS_URL"), paymentHash)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Api-Key", os.Getenv("LNBITS_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var payment struct {
		Paid bool `json:"paid"`
	}
	json.NewDecoder(resp.Body).Decode(&payment)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}
