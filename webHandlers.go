package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Handles requests to new.html
func pageNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")
	http.ServeFile(w, r, "templates/edit.html")
}

// Creates the page for https://host/view?id=userhash
func pageViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	http.ServeFile(w, r, "templates/view.html")
}

// Handles the json request to update a new meetup. If no adminhash is present, then a new meetup gets created.
func updateMeetUp(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Printf("updateMeetUp failed: %s\n", closeErr)
		}
	}()

	var newMeetUp MeetUp

	if err = readAndValidateJsonMeetUp(r, &newMeetUp); err != nil { // Dates and email address validated here
		writeJsonError(w, err.Error())
		return
	}

	for i := 0; i < len(newMeetUp.Dates); i++ {
		if len(newMeetUp.Dates[i].Users) > 0 { // No users allowed when creating or updating
			writeJsonError(w, "invalid user object.")
			return
		}
	}

	if newMeetUp.AdminHash == "" { // If no adminhash, a new meetup is being created, therefore generate both the hashes.
		// Generate the user hash
		randBytes := make([]byte, 64)
		if _, err := rand.Read(randBytes); err != nil {
			log.Printf("updateMeetUp failed: error reading random bytes for user hash. %s\n", err)
			writeJsonError(w, "Error reading random bytes.")
		}
		newMeetUp.UserHash = fmt.Sprintf("%x", sha256.Sum256(randBytes))

		// Generate the admin hash
		randBytes = make([]byte, 64)
		if _, err := rand.Read(randBytes); err != nil {
			log.Printf("updateMeetUp failed: error reading random bytes for admin hash. %s\n", err)
			writeJsonError(w, "Error reading random bytes.")
		}
		newMeetUp.AdminHash = fmt.Sprintf("%x", sha256.Sum256(randBytes))

		err = newMeetUp.CreateMeetUp()
		if err != nil {
			log.Printf("ajaxCreateHandler: err creating database rows: %s\n", err)
			writeJsonError(w, "Error creating new meetup.")
			return
		}
	} else {
		// Check the adminhash is valid
		if err = validateHash(newMeetUp.AdminHash); err != nil {
			log.Printf("updateMeetUp failed: invalid admin hash. hash:%q, error:%s\n", newMeetUp.AdminHash, err)
			writeJsonError(w, "invalid admin hash.")
			return
		}

		var currMeetUp MeetUp

		//Get MeetUp object by adminhash
		if err = currMeetUp.GetByAdminHash(newMeetUp.AdminHash); err != nil {
			if err.Error() == "no rows matching the adminhash" { // No rows found for this hash, send the user to the start page.
				log.Printf("updateMeetUp failed: hash:%q, error:%s\n", newMeetUp.AdminHash, err)
				writeJsonError(w, "admin hash not found.")
			} else {
				log.Printf("updateMeetUp failed: MeetUp.GetByAdminHash() hash:%q, error:%s\n", newMeetUp.AdminHash, err)
				writeJsonError(w, "database error.")
			}
			return
		}

		// Update the database.
		if err = currMeetUp.UpdateMeetUpDeleteDates(&newMeetUp); err != nil {
			log.Printf("updateMeetUp failed: MeetUp.UpdateMeetUpDeleteDates() hash:%q, error:%s\n", newMeetUp.AdminHash, err)
			writeJsonError(w, "database error. could not update.")
			return
		}

		newMeetUp.UserHash = currMeetUp.UserHash // hash missing in newMeetUp
	}

	// Create and write json response to the client
	type CreateResponseResult struct {
		UserHash  string `json:"userhash"`
		AdminHash string `json:"adminhash"`
	}
	type CreateResponse struct {
		Result CreateResponseResult `json:"result"`
		Error  string               `json:"error"`
	}

	successResponse := CreateResponse{Result: CreateResponseResult{UserHash: newMeetUp.UserHash, AdminHash: newMeetUp.AdminHash}, Error: ""}

	js, err := json.Marshal(successResponse)
	if err != nil {
		writeJsonError(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(js); err != nil {
		log.Printf("updateMeetUp failed: error writing response. %s\n", err)
	}
}

// Handles the json request to get meetup info with a user hash.
func getUserMeetUp(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	type reqStruct struct {
		UserHash string `json:"userhash"`
	}
	var reqJson reqStruct

	// Decode the json
	if err = json.NewDecoder(io.LimitReader(r.Body, 512)).Decode(&reqJson); err != nil { // 512B max json length
		log.Printf("getUserMeetUp invalid json: %s\n", err)
		writeJsonError(w, "invalid json.")
		return
	}

	// Check the userhash is valid
	if err = validateHash(reqJson.UserHash); err != nil {
		log.Printf("getUserMeetUp invalid user hash: %s\n", err)
		writeJsonError(w, "invalid hash.")
		return
	}

	meetUpObj := MeetUp{}

	if err = meetUpObj.GetByUserHash(reqJson.UserHash); err != nil {
		if err.Error() == "no rows matching the userhash" {
			writeJsonError(w, "user hash not found.")
		} else {
			log.Printf("getUserMeetUp: err getting by userhash: %s\n", err)
			writeJsonError(w, "database error.")
		}
		return
	}

	// Create and write json response to the client
	type CreateResponseResult struct {
		Dates       Dates  `json:"dates"`
		Description string `json:"description"`
	}
	type CreateResponse struct {
		Result CreateResponseResult `json:"result"`
		Error  string               `json:"error"`
	}

	successResponse := CreateResponse{Result: CreateResponseResult{Dates: meetUpObj.Dates, Description: meetUpObj.Description}, Error: ""}

	js, err := json.Marshal(successResponse)
	if err != nil {
		writeJsonError(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(js); err != nil {
		log.Printf("getUserMeetUp failed: error writing response. %s\n", err)
	}
}

// Handles the json request to get meetup info with an admin hash.
func getAdminMeetUp(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	type reqStruct struct {
		AdminHash string `json:"adminhash"`
	}
	var reqJson reqStruct

	// Decode the json
	if err = json.NewDecoder(io.LimitReader(r.Body, 512)).Decode(&reqJson); err != nil { // 512B max json length
		log.Printf("getAdminMeetUp invalid json: %s\n", err)
		writeJsonError(w, "invalid json.")
	}

	// Check the adminhash is valid
	if err = validateHash(reqJson.AdminHash); err != nil {
		log.Printf("getAdminMeetUp invalid admin hash: %s\n", err)
		writeJsonError(w, "invalid hash.")
		return
	}

	meetUpObj := MeetUp{}

	if err = meetUpObj.GetByAdminHash(reqJson.AdminHash); err != nil {
		if err.Error() == "no rows matching the adminhash" {
			writeJsonError(w, "admin hash not found.")
		} else {
			log.Printf("getAdminMeetUp: err getting by adminhash: %s\n", err)
			writeJsonError(w, "database error.")
		}
		return
	}

	// Create and write json response to the client
	type CreateResponse struct {
		Result MeetUp `json:"result"`
		Error  string `json:"error"`
	}

	successResponse := CreateResponse{Result: meetUpObj, Error: ""}

	js, err := json.Marshal(&successResponse)
	if err != nil {
		writeJsonError(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(js); err != nil {
		log.Printf("getAdminMeetUp failed: error writing response. %s\n", err)
	}
}

// Handles the json request to get meetup info with an admin hash.
func deleteMeetUp(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	type reqStruct struct {
		AdminHash string `json:"adminhash"`
	}
	var reqJson reqStruct

	// Decode the json
	if err = json.NewDecoder(io.LimitReader(r.Body, 512)).Decode(&reqJson); err != nil { // 512B max json length
		log.Printf("deleteMeetUp failed: invalid json: %s\n", err)
		writeJsonError(w, "invalid json.")
	}

	// Check the adminhash is valid
	if err = validateHash(reqJson.AdminHash); err != nil {
		log.Printf("deleteMeetUp failed: invalid admin hash: %s\n", err)
		writeJsonError(w, "invalid hash.")
		return
	}

	// Delete MeetUp object by adminhash
	var dbMeetUp MeetUp
	if err = dbMeetUp.DeleteByAdminHash(reqJson.AdminHash); err != nil {
		log.Printf("deleteMeetUp failed: err deleting MeetUp: %s\n", err)
		writeJsonError(w, "error deleting meetup")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	js := []byte(`{"result":"", "error":""}`)

	if _, err = w.Write(js); err != nil {
		log.Printf("deleteMeetUp failed: error writing response. %s\n", err)
	}
}

// Handles the json request to update a user. If the user is not present then they get added.
func updateUser(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	type reqStruct struct {
		UserHash string  `json:"userhash"`
		UserName string  `json:"username"`
		Dates    []int64 `json:"dates"`
	}
	var reqJson reqStruct

	// Decode the json
	if err = json.NewDecoder(io.LimitReader(r.Body, 4096)).Decode(&reqJson); err != nil { // 4KB max json length
		log.Printf("updateUser failed: invalid json: %s\n", err)
		writeJsonError(w, "invalid json.")
	}

	// Check the userhash is valid
	if err = validateHash(reqJson.UserHash); err != nil {
		log.Printf("updateUser failed: invalid user hash: %s\n", err)
		writeJsonError(w, "invalid hash.")
		return
	}

	meetUpObj := MeetUp{}

	if err = meetUpObj.GetByUserHash(reqJson.UserHash); err != nil {
		if err.Error() == "no rows matching the userhash" {
			writeJsonError(w, "user hash not found.")
		} else {
			log.Printf("updateUser: err getting by userhash: %s\n", err)
			writeJsonError(w, "database error.")
		}
		return
	}

	// Check if the user name is already in the database. All users are in each Date.Users slice, so only check the first one.
	var userPresent = false
	if len(meetUpObj.Dates) > 0 && len(meetUpObj.Dates[0].Users) > 0 {
		for _, userObj := range meetUpObj.Dates[0].Users {
			if userObj.Name == reqJson.UserName {
				userPresent = true
				break
			}
		}
	}

	if userPresent == true {
		users := Users{}

		for index, dateObj := range meetUpObj.Dates { // loop over dates
			for _, reqDate := range reqJson.Dates { // Get matching date to update. Only dates the user is available for are included in the query
				user := User{IdDate: meetUpObj.Dates[index].Id, Name: reqJson.UserName, Available: false}
				if reqDate == dateObj.Date {
					user.Available = true
				}
				users = append(users, user)
			}
		}

		if err = users.UpdateUsers(); err != nil {
			log.Printf("updateUser: error updating users: %s\n", err)
			writeJsonError(w, "database error updating user.")
			return
		}

	} else {
		// Go through database dates, try and get each date from the form, add user to each date, if date present in form will be "<timestamp>":"on" in form
		users := Users{}
		for index, dateObj := range meetUpObj.Dates {
			user := User{IdDate: meetUpObj.Dates[index].Id, Name: reqJson.UserName, Available: true}

			if len(r.FormValue(strconv.FormatInt(dateObj.Date, 10))) == 0 { // Not present, add user as unavailable
				user.Available = false
			}

			users = append(users, user)
		}

		if err = users.CreateUsers(); err != nil {
			log.Printf("updateUser: error creating users: %s\n", err)
			writeJsonError(w, "database error creating user.")
			return
		}
	}

	// Finished with the database return json
	w.Header().Set("Content-Type", "application/json")
	js := []byte(`{"result":"", "error":""}`)

	if _, err = w.Write(js); err != nil {
		log.Printf("updateUser, error writing response. %s\n", err)
	}
}

// Handles the json request to delete a user.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	type reqStruct struct {
		UserHash string `json:"userhash"`
		UserName string `json:"username"`
	}
	var reqJson reqStruct

	// Decode the json
	if err = json.NewDecoder(io.LimitReader(r.Body, 512)).Decode(&reqJson); err != nil { // 512B max json length
		log.Printf("deleteUser failed: invalid json: %s\n", err)
		writeJsonError(w, "invalid json.")
	}

	// Check the userhash is valid
	if err = validateHash(reqJson.UserHash); err != nil {
		log.Printf("deleteUser failed: invalid user hash: %s\n", err)
		writeJsonError(w, "invalid hash.")
		return
	}

	meetUpObj := MeetUp{}

	if err = meetUpObj.GetByUserHash(reqJson.UserHash); err != nil {
		if err.Error() == "no rows matching the userhash" {
			writeJsonError(w, "user hash not found.")
		} else {
			log.Printf("deleteUser: err getting by userhash: %s\n", err)
			writeJsonError(w, "database error.")
		}
		return
	}

	var userIds []int64

	for _, dateObj := range meetUpObj.Dates {
		for _, userObj := range dateObj.Users {
			if userObj.Name == reqJson.UserName {
				userIds = append(userIds, userObj.Id)
			}
		}
	}

	if len(userIds) > 0 {
		if err := DeleteByUserIds(userIds); err != nil {
			log.Printf("deleteUser: err deleting user: %s\n", err)
			writeJsonError(w, "database error.")
			return
		}
	}

	// Finished with the database return json
	w.Header().Set("Content-Type", "application/json")
	js := []byte(`{"result":"", "error":""}`)

	if _, err = w.Write(js); err != nil {
		log.Printf("updateUser, error writing response. %s\n", err)
	}
}
