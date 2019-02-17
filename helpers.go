package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/mail"
	"regexp"
)

// Extra helper functions that don't fit anywhere specifically

// Validates a hexadecimal sha256 hash
func validateHash(hash string) error {
	matched, err := regexp.MatchString("[^0-9a-fA-F]", hash)
	if err != nil {
		log.Printf("validateHash regexp.MatchString failed on hash %s: %v\n", hash, err)
		return err
	} else if matched == true {
		return errors.New("not hexadecimal")
	} else if len(hash) != 64 {
		return errors.New("not 64 bytes long")
	}
	return nil
}

// Returns a json error to the client
func writeJsonError(w http.ResponseWriter, errString string) {
	outMap := map[string]string{
		"result": "",
		"error":  errString,
	}
	js, err := json.Marshal(outMap)
	if err != nil {
		log.Printf("writeJsonError json marshalling failed: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(js); err != nil {
		log.Printf("writeJsonError failed writing the response: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Parses the json input into a MeetUp object.
// Returns an error suitable for returning to the user.
func readAndValidateJsonMeetUp(r *http.Request, meetUp *MeetUp) error {
	var err error

	// Decode the json into a MeetUp struct
	if err = json.NewDecoder(io.LimitReader(r.Body, 4096)).Decode(meetUp); err != nil { // 4KB max json length
		log.Printf("readAndValidateJsonMeetUp invalid json: %s\n", err)
		return errors.New("invalid json")
	}

	// Validate email address
	if _, err = mail.ParseAddress(meetUp.Admin.Email); err != nil {
		return errors.New("invalid email address")
	}

	// Validate dates
	if meetUp.Dates == nil || len(meetUp.Dates) == 0 {
		return errors.New("no dates selected")
	} else {
		for i := 0; i < len(meetUp.Dates); i++ {
			if meetUp.Dates[i].Date <= 0 {
				return errors.New("invalid date")
			}
		}
	}
	return nil
}
