package json

import (
	"encoding/json"
	"io"
	"net/http"
)

func ConvertTOJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

// DecodeFromJSON decodes the JSON request body into result.
// Returns nil when the body is empty (io.EOF) — callers that require
// a body should enforce that via struct validation tags.
func DecodeFromJSON(r *http.Request, result any) error {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(result); err != nil {
		if err == io.EOF {
			return nil // empty body is not an error
		}
		return err
	}
	return nil
}