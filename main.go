package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	fmt.Println("hello world")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		blockheight := "1111111"
		tmpl.Execute(w, blockheight)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		// log.Print("HTMX request recieved")
		// log.Print(r.Header.Get("HX-Request"))
		year := r.PostFormValue("year")
		fmt.Println(year)
	}
	http.HandleFunc("/", h1)
	http.HandleFunc("/get-blockheight/", h2)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
