package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type jsonUpload struct {
	uploadTime time.Time
	contents   []byte
}

var jsonUploads []jsonUpload

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"raw": func(ju jsonUpload) template.HTML {
		return template.HTML(string(ju.contents))
	},
}).ParseGlob("*.html"))

func main() {
	jsonUploads = []jsonUpload{}

	http.HandleFunc("/", serveJSON)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveJSON(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to ", r.URL.Path)

	if r.URL.Path != "/" {
		s := strings.Split(r.URL.Path, "/")[1]
		i, err := strconv.Atoi(s)

		if i >= len(jsonUploads) || err != nil {
			return
		}

		ju := jsonUploads[i]
		rd := bytes.NewReader(ju.contents)
		http.ServeContent(w, r, fmt.Sprintf("file%d.json", i), ju.uploadTime, rd)
		return
	}

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
			jsonUploads = append(jsonUploads, jsonUpload{time.Now(), b})
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
