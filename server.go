package pharmafinder

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"pharmafinder/utils"
	"regexp"

	"github.com/rs/zerolog"
)

const PATH_PREFIX = "frontend/build"

var logger zerolog.Logger = utils.GetLogger("WEB")

// Static file handler
func StaticServer(w http.ResponseWriter, r *http.Request) {
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
		logger.Error().Msgf("Could not open file %s: %v\n", r.URL.Path, err)
		logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("addr", r.RemoteAddr).
			Int("code", http.StatusInternalServerError).
			Msg("Request made")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error().Msgf("Failed to read file %s: %v\n", r.URL.Path, err)
		logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("addr", r.RemoteAddr).
			Int("code", http.StatusInternalServerError).
			Msg("Request made")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", mimetype)
	logger.Debug().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("addr", r.RemoteAddr).
		Int("code", http.StatusOK).
		Msg("Request made")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
