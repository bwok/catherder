package main

import (
	"log"
	"net/http"
)

// Routes all non /api/... requests
func defaultRouter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")

	switch r.URL.Path {
	case "/edit":
		pageEditHandler(w, r)
		break
	case "/view":
		pageViewHandler(w, r)
		break
	default:
		t, httpCode := templateJobber("/index", nil)
		if httpCode > 0 {
			http.Error(w, http.StatusText(httpCode), httpCode)
			if closeErr := r.Body.Close(); closeErr != nil {
				log.Println(closeErr)
			}
			return
		} else if err := t.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
			log.Println(err)
		}
	}
}

// Handles requests to new.html
func pageEditHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")

	t, httpCode := templateJobber(r.URL.Path, nil)
	if httpCode > 0 {
		http.Error(w, http.StatusText(httpCode), httpCode)
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
		return
	}

	data := struct {
		Title  string
		Header string
	}{
		Title:  "Create a meet up",
		Header: "Create Your Meet Up",
	}

	if err := validateHash(r.FormValue("id")); err == nil { // On valid adminhash, change the page title
		data.Title = "Edit your meet up"
		data.Header = "Edit Your Meet Up"
	}

	err := t.ExecuteTemplate(w, "edit.gohtml", data)
	if err != nil {
		log.Println(err)
	}
}

// Creates the page for https://host/view?id=userhash
func pageViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	t, httpCode := templateJobber(r.URL.Path, nil)
	if httpCode > 0 {
		http.Error(w, http.StatusText(httpCode), httpCode)
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
		return
	}
	err := t.ExecuteTemplate(w, "view.gohtml", nil)
	if err != nil {
		log.Println(err)
	}
}
