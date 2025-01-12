package main

import (
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	blockheight := "Enter a date to find the blockheight."
	tmpl.Execute(w, blockheight)
}
