package pharmafinder

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"regexp"
)

const PATH_PREFIX = "frontend/build"

// Static file handler
func StaticServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request made to %s\n", r.URL.Path)
	regex := regexp.MustCompile(`^.*/(.*?(\.[A-Za-z0-9]+))$`)

	var path string
	if regex.MatchString(r.URL.Path) {
		path = fmt.Sprintf("%s%s", PATH_PREFIX, r.URL.Path)
	} else if r.URL.Path[len(r.URL.Path)-1] != '/' {
		path = fmt.Sprintf("%s%s.html", PATH_PREFIX, r.URL.Path)
	} else {
		path = fmt.Sprintf("%s%sindex.html", PATH_PREFIX, r.URL.Path)
	}

	file, err := ServerFS.Open(path)
	ext := regex.FindStringSubmatch(path)[2]
	mimetype := mime.TypeByExtension(ext)

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
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
