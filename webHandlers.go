package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
			UserLink  string `json:"userlink"`
			AdminLink string `json:"adminlink"`
		}
		type CreateResponse struct {
			Result CreateResponseResult `json:"result"`
			Error  string               `json:"error"`
		}

		successResponse := CreateResponse{Result: CreateResponseResult{UserLink: newMeetUp.UserHash, AdminLink: newMeetUp.AdminHash}, Error: ""}

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

// Creates the page for https://host/view?id=userhash
func pageViewHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; img-src 'self'; style-src 'self';")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	userHash := r.FormValue("id")
	meetUpObj := MeetUp{}

	// Check the userhash is valid
	if err = validateHash(userHash); err != nil {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	if err = meetUpObj.GetByUserHash(userHash); err != nil {
		if err.Error() == "no rows matching the userhash" { // No rows found for this hash, send the user to the start page.
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusFound)
		} else {
			log.Printf("pageViewHandler: err getting by userhash: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var viewTemplate = template.Must(template.New("view.html").
		Funcs(template.FuncMap{"getMonth": getMonth, "getDate": getDate, "getWeekDay": getWeekDay}).
		ParseFiles("templates/view.html"))

	type viewObject struct {
		Host string
		MeetUp
	}

	view := viewObject{Host: r.Host, MeetUp: meetUpObj}

	if err = viewTemplate.Execute(w, view); err != nil {
		log.Printf("executing template failed: %s\n", err)
	}
}

// Called when adding a new user
func addUserHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	if err := r.ParseForm(); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userName := r.FormValue("username")
	if len(userName) == 0 {
		log.Println("Username either absent or with length 0 in /adduser request.")
		http.Error(w, "Invalid username.", http.StatusBadRequest)
		return
	}

	userHash := r.FormValue("userhash")
	if len(userHash) == 0 {
		log.Println("userHash either absent or with length 0 in /adduser request.")
		http.Error(w, "Invalid url.", http.StatusBadRequest)
		return
	}

	meetUpObj := MeetUp{}

	if err = meetUpObj.GetByUserHash(userHash); err != nil {
		if err.Error() == "no rows" { // No rows found for this hash, send the user to the start page.
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusFound)
		} else {
			log.Printf("addUserHandler: err getting by userhash: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Go through database dates, try and get each date from the form, add user to each date, if date present in form will be "<timestamp>":"on" in form
	users := Users{}
	for index, dateObj := range meetUpObj.Dates {
		user := User{IdDate: meetUpObj.Dates[index].Id, Name: userName}

		if len(r.FormValue(strconv.FormatInt(dateObj.Date, 10))) == 0 { // Not present, add user as unavailable
			user.Available = false
		} else {
			user.Available = true
		}

		users = append(users, user)
	}

	if err = users.CreateUsers(); err != nil {
		log.Printf("addUserHandler: error creating users: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/view?id="+userHash)
	w.WriteHeader(http.StatusFound)
}

// Handles requests to admin.html
func pageAdminHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	if err := r.ParseForm(); err != nil {
		log.Printf("pageAdminHandler: error parsing form %s.\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	adminHash := r.FormValue("id")
	meetUpObj := MeetUp{}

	// Check the adminhash is valid
	if err = validateHash(adminHash); err != nil {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	if err = meetUpObj.GetByAdminHash(adminHash); err != nil {
		if err.Error() == "no rows matching the adminhash" { // No rows found for this hash, send the user to the start page.
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusFound)
		} else {
			log.Printf("pageAdminHandler: err getting by adminhash: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var adminTemplate = template.Must(template.New("admin.html").
		Funcs(template.FuncMap{"getMonth": getMonth, "getDate": getDate, "getWeekDay": getWeekDay}).
		ParseFiles("templates/admin.html"))

	if err = adminTemplate.Execute(w, meetUpObj); err != nil {
		log.Printf("executing template failed: %s\n", err)
	}
}

// Updates an edited MeetUp.
func ajaxAdminSaveHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Printf("ajaxAdminSaveHandler failed: %s\n", closeErr)
		}
	}()

	adminHash := r.FormValue("id")
	var newMeetUp MeetUp
	var currMeetUp MeetUp

	if err = readAndValidateJsonMeetUp(r, &newMeetUp); err != nil {
		writeJsonError(w, err.Error())
		return
	}

	// Check the adminhash is valid
	if err = validateHash(adminHash); err != nil {
		http.Error(w, "Invalid URL", http.StatusInternalServerError)
		return
	}

	//Get MeetUp object by adminhash
	if err = currMeetUp.GetByAdminHash(adminHash); err != nil {
		if err.Error() == "no rows matching the adminhash" { // No rows found for this hash, send the user to the start page.
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusFound)
		} else {
			log.Printf("ajaxAdminSaveHandler: err getting by adminhash: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Update the database.
	if err = currMeetUp.UpdateMeetUpDeleteDates(&newMeetUp); err != nil {
		log.Printf("ajaxAdminSaveHandler: err updating MeetUp: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type Response struct {
		Result string `json:"result"`
		Error  string `json:"error"`
	}

	successResponse := Response{Result: currMeetUp.UserHash, Error: ""}

	js, err := json.Marshal(successResponse)
	if err != nil {
		writeJsonError(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(js); err != nil {
		log.Printf("ajaxAdminSaveHandler, error writing response. %s\n", err)
	}
}

// Deletes a MeetUp.
func ajaxAdminDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Printf("ajaxAdminDeleteHandler failed: %s\n", closeErr)
		}
	}()

	var dbMeetUp MeetUp
	adminHash := r.FormValue("id")

	// Delete MeetUp object by adminhash
	if err = dbMeetUp.DeleteByAdminHash(adminHash); err != nil {
		log.Printf("ajaxAdminDeleteHandler: err deleting MeetUp: %s\n", err)
		writeJsonError(w, "error deleting meetup")
		return
	}

	type Response struct {
		Result string `json:"result"`
		Error  string `json:"error"`
	}

	successResponse := Response{Result: "", Error: ""}

	js, err := json.Marshal(successResponse)
	if err != nil {
		writeJsonError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(js); err != nil {
		log.Printf("ajaxAdminDeleteHandler, error writing response. %s\n", err)
	}

}
