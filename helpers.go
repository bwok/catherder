package main

import (
	"encoding/binary"
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
	if _, err = mail.ParseAddress(meetUp.AdminEmail); err != nil {
		return errors.New("invalid email address")
	}

	// Validate dates
	if len(meetUp.Dates) == 0 {
		return errors.New("no dates selected")
	} else {
		for _, date := range meetUp.Dates {
			if date <= 0 {
				return errors.New("invalid date")
			}
		}
	}
	return nil
}

// convert dates []int64 to a []byte (sqlite blob type).
func convertDatesToBlob(intSlice []int64) []byte {
	var blobBytes []byte

	for _, date := range intSlice {
		buf := make([]byte, binary.MaxVarintLen64)
		binary.PutVarint(buf, date)
		blobBytes = append(blobBytes, buf...)
	}

	return blobBytes
}

// convert blob []byte to date []int64
func convertBlobToDates(blobBytes []byte) []int64 {
	var intSlice = make([]int64, 0)

	for i := 0; i < len(blobBytes); i += binary.MaxVarintLen64 {
		outInt, _ := binary.Varint(blobBytes[i : i+binary.MaxVarintLen64])
		intSlice = append(intSlice, outInt)
	}

	return intSlice
}
