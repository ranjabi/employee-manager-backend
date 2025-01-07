package lib

import "net/http"

func WriteJsonResponse(w http.ResponseWriter, statusCode int, res []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(res)
}
