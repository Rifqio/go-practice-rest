package lib

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(res http.ResponseWriter, status int, data any) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(data)
}
