package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var jsonUploads []template.HTML

var templates = template.Must(template.ParseGlob("*.html"))

func main() {
	jsonUploads = []template.HTML{}

	http.HandleFunc("/", serveJSON)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveJSON(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		uploadJSON(w, r)
	case "GET":
		showJSON(w, r)
	}
}

func uploadJSON(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var u interface{}

	if err := decoder.Decode(&u); err == nil {
		if b, err := json.MarshalIndent(u, "", " "); err == nil {
			jsonUploads = append(jsonUploads, template.HTML(string(b)))
		} else {
			log.Println("Error indenting json: ", err)
		}
	} else {
		log.Println("Error parsing json: ", err)
	}
}

func showJSON(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "index.html", jsonUploads); err != nil {
		log.Println("Error rendering template: ", err)
	}
}
