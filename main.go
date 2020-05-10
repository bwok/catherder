package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

func main() {
	var err error

	connStr := "dbname=meetupdatabase user=meetupuser password=testpassword host=192.168.56.51 "
	if db, err = sql.Open("postgres", connStr); err != nil {
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
	httpsServeMux.HandleFunc("/api/", apiRouter)			// JSON request/response handlers
	httpsServeMux.HandleFunc("/", defaultRouter)			// All non /served/ or /api/ requests

	httpsServ := &http.Server{
		Addr: ":" + strconv.Itoa(ConfigObj.HttpsPort),
	}
	httpsServ.Handler = httpsServeMux
	log.Fatal(httpsServ.ListenAndServeTLS("./cert.pem", "./key.pem"))
}
