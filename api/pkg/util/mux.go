package util

import (
	"net/http"
	"strconv"
)

func WriteJSONResponse(w http.ResponseWriter, status int, body []byte) {
	var bl string
	if body == nil {
		bl = "0"
	} else {
		bl = strconv.Itoa(len(body))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", bl)
	w.WriteHeader(status)
	w.Write(body)
}
