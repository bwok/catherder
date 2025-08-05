package main

import (
	"database/sql"
	"embed"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

const maxLongJsonBytesLen = 4096 // Limit create/update JSON requests to this many bytes
const maxShortJsonBytesLen = 512 // Limit the other JSON requests to this many bytes

var (
	//go:embed all:served
	served embed.FS

	//go:embed templates
	res   embed.FS
	pages = map[string]string{
		"/index": "templates/index.gohtml",
		"/edit":  "templates/edit.gohtml",
		"/view":  "templates/view.gohtml",
	}
)

func main() {
	port := flag.String("port", "443", "-port=<port> The port to listen for https requests on.")
	certPath := flag.String("cert", "./cert.pem", "-cert=<path> The path of the ssl certificate.")
	keyPath := flag.String("key", "./key.pem", "-key=<path> The path of the ssl key.")
	flag.Parse()

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
			if _, err = db.Exec(`
			PRAGMA foreign_keys = ON;

			CREATE TABLE IF NOT EXISTS meetup
			(
				idmeetup    INTEGER PRIMARY KEY ASC,
				userhash    TEXT    NOT NULL,
				adminhash   TEXT    NOT NULL,
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
			`); err != nil {
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

	// Serve https traffic
	servedDir, err := fs.Sub(served, "served")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/served/", http.StripPrefix("/served/", http.FileServer(http.FS(servedDir))))
	http.HandleFunc("/api/", apiRouter) // JSON request/response handlers
	http.HandleFunc("/", defaultRouter) // All non /served/ or /api/ requests
	log.Printf("Server starting up, listening at: :%s", *port)
	log.Fatal(http.ListenAndServeTLS(":"+*port, *certPath, *keyPath, nil))
}

// templateJobber Parses a cached template file
func templateJobber(path string, funcMap *template.FuncMap) (*template.Template, int) {
	filePath, ok := pages[path]
	if !ok {
		log.Printf("filePath %s not found in pages cache...", path)
		return nil, http.StatusNotFound
	}

	t := template.New("")
	if funcMap != nil {
		t.Funcs(*funcMap)
	}

	t, err := t.ParseFS(res, filePath)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		return nil, http.StatusInternalServerError
	}
	return t, -1
}
