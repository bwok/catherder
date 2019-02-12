package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Web server handlers

// Handles requests to new.html
func pageNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")
	http.ServeFile(w, r, "templates/new.html")
}

// Handles the ajax request to create a new meetup
func ajaxCreateHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Printf("ajaxCreateHandler failed: %s\n", closeErr)
		}
	}()

	var newMeetUp MeetUp

	if err = readAndValidateJsonMeetUp(r, &newMeetUp); err != nil {
		writeJsonError(w, err.Error())
		return
	}

	for i := 0; i < len(newMeetUp.Dates); i++ {
		if len(newMeetUp.Dates[i].Users) > 0 { // No users allowed when creating
			writeJsonError(w, "invalid user object.")
			return
		}
	}

	// Generate the user hash
	randBytes := make([]byte, 64)
	if _, err := rand.Read(randBytes); err != nil {
		log.Printf("ajaxCreateHandler failed: error reading random bytes for user hash. %s\n", err)
		writeJsonError(w, "Error reading random bytes.")
	}
	newMeetUp.UserHash = fmt.Sprintf("%x", sha256.Sum256(randBytes))

	// Generate the admin hash
	randBytes = make([]byte, 64)
	if _, err := rand.Read(randBytes); err != nil {
		log.Printf("ajaxCreateHandler failed: error reading random bytes for admin hash. %s\n", err)
		writeJsonError(w, "Error reading random bytes.")
	}
	newMeetUp.AdminHash = fmt.Sprintf("%x", sha256.Sum256(randBytes))

	// Create the new meetup rows in the database.
	err = newMeetUp.CreateMeetUp()
	if err != nil {
		log.Printf("ajaxCreateHandler: err creating database rows: %s\n", err)
		writeJsonError(w, "Error creating new meetup.")
	} else {
		type CreateResponseResult struct {
			UserLink string `json:"userlink"`
			EditLink string `json:"editlink"`
		}
		type CreateResponse struct {
			Result CreateResponseResult `json:"result"`
			Error  string               `json:"error"`
		}

		successResponse := CreateResponse{Result: CreateResponseResult{UserLink: newMeetUp.UserHash, EditLink: newMeetUp.AdminHash}, Error: ""}

		js, err := json.Marshal(successResponse)
		if err != nil {
			writeJsonError(w, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(js)
		if err != nil {
			log.Printf("ajaxCreateHandler, error writing response. %s\n", err)
		}
	}
}
