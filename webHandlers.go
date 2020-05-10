package main

import (
	"html/template"
	"log"
	"net/http"
)

const maxLongJsonBytesLen = 4096 // Limit create/update JSON requests to this many bytes
const maxShortJsonBytesLen = 512 // Limit the other JSON requests to this many bytes

// Routes all non /api/... requests
func defaultRouter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")

	switch r.URL.Path {
	case "/edit":
		pageEditHandler(w, r)
	case "/view":
		pageViewHandler(w, r)
	default:
		http.ServeFile(w, r, "templates/index.html")
	}
}

// Handles requests to new.html
func pageEditHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")

	t := template.Must(template.ParseFiles("templates/edit.html"))

	data := struct {
		Title string
	}{
		Title: "Create a meet up",
	}

	if err := validateHash(r.FormValue("id")); err == nil { // On valid adminhash, change the page title
		data.Title = "Edit your meet up"
	}

	err := t.ExecuteTemplate(w, "edit", data)
	if err != nil {
		log.Fatal(err)
	}
}

// Creates the page for https://host/view?id=userhash
func pageViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	http.ServeFile(w, r, "templates/view.html")
}
