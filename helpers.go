package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// Extra helper functions that don't fit anywhere specifically

// Validates a hexadecimal sha256 hash
func validateHash(hash string) error {
	hashByteLen := 64

	matched, err := regexp.MatchString("[^0-9a-fA-F]", hash)
	if err != nil {
		log.Printf("validateHash regexp.MatchString failed on hash %s: %v\n", hash, err)
		return err
	} else if matched == true {
		return errors.New("not hexadecimal")
	} else if len(hash) != hashByteLen {
		return errors.New("not " + strconv.Itoa(hashByteLen) + " bytes long")
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
