package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strconv"
)

var db *sql.DB

func main() {
	var err error

	// Check the database file is present, open it
	if _, err = os.Stat("./data.sqlite"); err == nil {
		if db, err = sql.Open("sqlite3", "file:data.sqlite?_foreign_keys=true"); err != nil {
			log.Fatal(err)
		}
	} else if os.IsNotExist(err) {
		//log.Fatal(err)
		if db, err = sql.Open("sqlite3", "file:data.sqlite?_foreign_keys=true"); err != nil {
			log.Fatal(err)
		} else {
			sqlStmt := `
			CREATE TABLE IF NOT EXISTS meetup
			(
				idmeetup    INTEGER PRIMARY KEY ASC,
				userhash    TEXT    NOT NULL,
				adminhash   TEXT    NOT NULL,
				adminemail  TEXT    NOT NULL,
				sendalerts  INTEGER NOT NULL,
				dates       BLOB    NOT NULL,
				description TEXT    NOT NULL
			);
			
			CREATE TABLE IF NOT EXISTS "user"
			(
				iduser   INTEGER PRIMARY KEY ASC NOT NULL,
				idmeetup INTEGER                 NOT NULL,
				name     TEXT                    NOT NULL,
				dates    BLOB                    NOT NULL,
				FOREIGN KEY (idmeetup) REFERENCES meetup (idmeetup) ON DELETE CASCADE
			);
			
			CREATE INDEX IF NOT EXISTS "user.fk_user_meetup_idx" ON "user" ("idmeetup");
			`
			if _, err = db.Exec(sqlStmt); err != nil {
				log.Fatal(err)
			}
		}
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	// Prepare all the sql statements for later use
	prepareDatabaseStatements()
	defer closeDatabaseStatements()

	// Redirect http traffic to https
	httpServ := &http.Server{
		Addr: ":" + strconv.Itoa(ConfigObj.HttpPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Connection", "close")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			url := "https://" + r.Host + r.URL.String()
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}),
	}
	go func() { log.Fatal(httpServ.ListenAndServe()) }()

	// Serve https traffic
	httpsServeMux := http.NewServeMux()
	httpsServeMux.Handle("/served/", http.StripPrefix("/served/", http.FileServer(http.Dir("./served/"))))
	httpsServeMux.HandleFunc("/api/", apiRouter) // JSON request/response handlers
	httpsServeMux.HandleFunc("/", defaultRouter) // All non /served/ or /api/ requests

	httpsServ := &http.Server{
		Addr: ":" + strconv.Itoa(ConfigObj.HttpsPort),
	}
	httpsServ.Handler = httpsServeMux
	log.Fatal(httpsServ.ListenAndServeTLS("./cert.pem", "./key.pem"))
}
