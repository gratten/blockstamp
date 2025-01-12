package main

import (
	"html/template"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Load the layout and the home content
	tmpl, err := template.ParseFiles(
		"app/layout.html",
		"app/home.html",
		"app/blockheight.html",
	)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	// Dynamic data (you can replace "blockheight" with actual logic)
	blockheight := "Enter a date to find the blockheight."

	// Execute the template, passing the dynamic data for blockheight
	err = tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Title":       "Home",
		"Blockheight": blockheight,
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render home page", http.StatusInternalServerError)
		return
	}
}
