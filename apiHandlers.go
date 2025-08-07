package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Routes all /api/... requests
func apiRouter(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/updatemeetup":
		updateMeetUp(w, r)
		break
	case "/api/getusermeetup":
		getUserMeetUp(w, r)
		break
	case "/api/getadminmeetup":
		getAdminMeetUp(w, r)
		break
	case "/api/deletemeetup":
		deleteMeetUp(w, r)
		break
	case "/api/updateuser":
		updateUser(w, r)
		break
	case "/api/deleteuser":
		deleteUser(w, r)
		break
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
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

	// Decode the json into a MeetUp struct
	if err = json.NewDecoder(io.LimitReader(r.Body, maxLongJsonBytesLen)).Decode(&newMeetUp); err != nil { // 4KB max json length
		log.Printf("updateMeetUp invalid json: %s\n", err)
		writeJsonError(w, "invalid json")
		return
	}

	// Validate dates
	if len(newMeetUp.Dates) == 0 {
		writeJsonError(w, "no dates selected")
		return
	} else {
		for _, date := range newMeetUp.Dates {
			if date <= 0 {
				writeJsonError(w, "invalid date")
				return
			}
		}
	}

	for i := 0; i < len(newMeetUp.Dates); i++ {
		if len(newMeetUp.Users) > 0 { // No users allowed when creating or updating
			writeJsonError(w, "invalid user object.")
			return
		}
	}

	if newMeetUp.AdminHash == "" { // If no adminhash, a new meetup is being created, therefore generate both the hashes.
		randByteLen := 64

		// Generate the user hash
		randBytes := make([]byte, randByteLen)
		if _, err := rand.Read(randBytes); err != nil {
			log.Printf("updateMeetUp failed: error reading random bytes for user hash. %s\n", err)
			writeJsonError(w, "Error reading random bytes.")
		}
		newMeetUp.UserHash = fmt.Sprintf("%x", sha512.Sum512(randBytes))

		// Generate the admin hash
		randBytes = make([]byte, randByteLen)
		if _, err := rand.Read(randBytes); err != nil {
			log.Printf("updateMeetUp failed: error reading random bytes for admin hash. %s\n", err)
			writeJsonError(w, "Error reading random bytes.")
		}
		newMeetUp.AdminHash = fmt.Sprintf("%x", sha512.Sum512(randBytes))

		err = newMeetUp.Create()
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
		currMeetUp.Dates = newMeetUp.Dates
		currMeetUp.Description = newMeetUp.Description

		if err = currMeetUp.Update(); err != nil {
			log.Printf("updateMeetUp failed: MeetUp.Update() hash:%q, error:%s\n", newMeetUp.AdminHash, err)
			writeJsonError(w, "database error. could not update.")
			return
		}

		newMeetUp = currMeetUp
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
	if err = json.NewDecoder(io.LimitReader(r.Body, maxShortJsonBytesLen)).Decode(&reqJson); err != nil {
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
			writeJsonError(w, "The meetup was not found.")
		} else {
			log.Printf("getUserMeetUp: err getting by userhash: %s\n", err)
			writeJsonError(w, "database error.")
		}
		return
	}

	// Create and write json response to the client
	type CreateResponseResult struct {
		Dates       []int64 `json:"dates"`
		Users       Users   `json:"users"`
		Description string  `json:"description"`
	}
	type CreateResponse struct {
		Result CreateResponseResult `json:"result"`
		Error  string               `json:"error"`
	}

	successResponse := CreateResponse{Result: CreateResponseResult{Dates: meetUpObj.Dates, Users: meetUpObj.Users, Description: meetUpObj.Description}, Error: ""}

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
	if err = json.NewDecoder(io.LimitReader(r.Body, maxShortJsonBytesLen)).Decode(&reqJson); err != nil {
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
			writeJsonError(w, "The meetup was not found.")
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
	if err = json.NewDecoder(io.LimitReader(r.Body, maxShortJsonBytesLen)).Decode(&reqJson); err != nil {
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
	if err = json.NewDecoder(io.LimitReader(r.Body, maxLongJsonBytesLen)).Decode(&reqJson); err != nil {
		log.Printf("updateUser failed: invalid json: %s\n", err)
		writeJsonError(w, "invalid json.")
	}

	// Check the userhash is valid
	if err = validateHash(reqJson.UserHash); err != nil {
		log.Printf("updateUser failed: invalid user hash: %s\n", err)
		writeJsonError(w, "invalid hash.")
		return
	}

	// Check the username is not empty
	if reqJson.UserName == "" {
		writeJsonError(w, "The user name is empty.")
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

	// Try and update an existing user with the same name, if the user is already in the database.
	var userPresent = false
	for _, userObj := range meetUpObj.Users {
		if userObj.Name == reqJson.UserName {
			userObj.Dates = reqJson.Dates
			if err = userObj.Update(); err != nil {
				log.Printf("updateUser: error updating users: %s\n", err)
				writeJsonError(w, "database error updating user.")
				return
			}
			userPresent = true
			break
		}
	}
	// Update failed, create a new user
	if userPresent == false {
		user := User{IdMeetUp: meetUpObj.Id, Name: reqJson.UserName, Dates: reqJson.Dates}
		if err = user.Create(); err != nil {
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
	if err = json.NewDecoder(io.LimitReader(r.Body, maxShortJsonBytesLen)).Decode(&reqJson); err != nil {
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

	for _, userObj := range meetUpObj.Users {
		if userObj.Name == reqJson.UserName {
			if err := userObj.Delete(); err != nil {
				log.Printf("deleteUser: err deleting user: %s\n", err)
				writeJsonError(w, "database error.")
				return
			}
		}
	}

	// Finished with the database return json
	w.Header().Set("Content-Type", "application/json")
	js := []byte(`{"result":"", "error":""}`)

	if _, err = w.Write(js); err != nil {
		log.Printf("updateUser, error writing response. %s\n", err)
	}
}
