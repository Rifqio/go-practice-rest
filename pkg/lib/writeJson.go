package lib

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(res http.ResponseWriter, status int, data any) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Set-Cookie", "sessionid=38afes7a8; HttpOnly; Secure")
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(data)
}
