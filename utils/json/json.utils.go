package json

import (
	"encoding/json"
	"net/http"
)

func ConvertTOJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func DecodeFromJSON(r *http.Request, result any) error {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(result)
}