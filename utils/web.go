package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, code int, resp interface{}) {
	b, _ := json.Marshal(resp)
	w.WriteHeader(code)
	w.Write(b)
}
