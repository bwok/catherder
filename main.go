package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var db *sql.DB

func main() {
	var err error

	// Check the database file is present, open it
	if _, err = os.Stat("./data.sqlite"); err == nil {
		if db, err = sql.Open("sqlite3", "./data.sqlite"); err != nil {
			log.Fatal(err)
		}
	} else if os.IsNotExist(err) {
		log.Fatal(err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	// Prepare all the sql statements for later use
	prepareDatabaseStatements()
	defer closeDatabaseStatements()

}
