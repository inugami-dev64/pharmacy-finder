package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"pharmafinder"
	"pharmafinder/db"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const PATH_PREFIX = "frontend/build"

func StaticServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request made to %s\n", r.URL.Path)
	regex := regexp.MustCompile(`^.*/(.*?(\.[A-Za-z0-9]+))$`)

	var file fs.File
	var err error
	var mimetype string
	if regex.MatchString(r.URL.Path) {
		file, err = pharmafinder.ServerFS.Open(fmt.Sprintf("%s%s", PATH_PREFIX, r.URL.Path))
		ext := regex.FindStringSubmatch(r.URL.Path)[2]
		mimetype = mime.TypeByExtension(ext)
	} else if r.URL.Path[len(r.URL.Path)-1] != '/' {
		path := fmt.Sprintf("%s.html", r.URL.Path)
		file, err = pharmafinder.ServerFS.Open(fmt.Sprintf("%s%s", PATH_PREFIX, path))
		ext := regex.FindStringSubmatch(path)[2]
		mimetype = mime.TypeByExtension(ext)
	} else {
		path := fmt.Sprintf("%sindex.html", r.URL.Path)
		file, err = pharmafinder.ServerFS.Open(fmt.Sprintf("%s%s", PATH_PREFIX, path))
		ext := regex.FindStringSubmatch(path)[2]
		mimetype = mime.TypeByExtension(ext)
	}

	if err != nil {
		log.Printf("Could not open file %s: %v\n", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Failed to read file %s: %v\n", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", mimetype)
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Attempt to load .env files if they exist
	godotenv.Load("deploy/.env")

	// Connect to the database
	conn := db.ConnectToDB()
	db.EnsureMigrationsAreUpToDate(conn)

	r := mux.NewRouter()
	r.PathPrefix("/").
		Methods("GET").
		HandlerFunc(StaticServer)

	server := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
