package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
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

	// Redirect http traffic to https
	httpServ := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Connection", "close")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			url := "https://" + r.Host + r.URL.String()
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	httpServ.SetKeepAlivesEnabled(false)
	go func() { log.Fatal(httpServ.ListenAndServe()) }()

	// Serve https traffic
	httpsServeMux := http.NewServeMux()
	httpsServeMux.Handle("/served/", http.StripPrefix("/served/", http.FileServer(http.Dir("./served/"))))
	httpsServeMux.HandleFunc("/edit", pageEditHandler)
	httpsServeMux.HandleFunc("/view", pageViewHandler)

	// JSON handlers
	httpsServeMux.HandleFunc("/api/updatemeetup", updateMeetUp)
	httpsServeMux.HandleFunc("/api/getusermeetup", getUserMeetUp)
	httpsServeMux.HandleFunc("/api/getadminmeetup", getAdminMeetUp)
	httpsServeMux.HandleFunc("/api/deletemeetup", deleteMeetUp)
	httpsServeMux.HandleFunc("/api/updateuser", updateUser)
	httpsServeMux.HandleFunc("/api/deleteuser", deleteUser)

	httpsServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';")
		http.ServeFile(w, r, "templates/index.html")
	})

	httpsServ := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	httpsServ.SetKeepAlivesEnabled(false)
	httpsServ.Handler = httpsServeMux
	log.Fatal(httpsServ.ListenAndServeTLS("./cert.pem", "./key.pem"))
}
